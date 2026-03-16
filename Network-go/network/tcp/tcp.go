package tcp

import (
	"bufio"
	"encoding/json"

	"heisprosjekt75/types"
	"log"
	"net"
	"strings"
	"time"
)

func readLoop(conn net.Conn, incomingTCP chan Message) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n') //hvor langt skal en string være? dette må skrives inn i parantesen
		if err != nil {
			log.Println("Message reading error in read loop: ", err)
			return
		}

		line = strings.TrimSpace(line)
		var msg Message
		err = json.Unmarshal([]byte(line), &msg)
		if err != nil {
			log.Println("json decode error in read loop: ", err)
			continue
		}
		incomingTCP <- msg
	}
}

func StartPrimaryTCP(ps *types.PeerState, port string, incomingTCP chan Message, e *types.Elevator) {

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println("Error with listning object")
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Error with listning conn")
				continue
			}
			go handleNewNode(conn, incomingTCP, e, ps)
		}
	}()
}

func handleNewNode(conn net.Conn, incomingTCP chan Message, e *types.Elevator, ps *types.PeerState) {
	reader := bufio.NewReader(conn)

	line, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Message reading error in handleNewNode:", err)
		_ = conn.Close()
		return
	}

	line = strings.TrimSpace(line)

	var msg Message
	err = json.Unmarshal([]byte(line), &msg)
	if err != nil {
		log.Println("json decode error in handleNewNode:", err)
		_ = conn.Close()
		return
	}

	if msg.Type != Msghello {
		log.Println("expected Msghello, got", msg.Type)
		_ = conn.Close()
		return
	}

	var hello HelloMessage
	helloBytes, err := json.Marshal(msg.MessageData)
	if err != nil {
		log.Println("marshal hello data failed:", err)
		_ = conn.Close()
		return
	}

	if err := json.Unmarshal(helloBytes, &hello); err != nil {
		log.Println("unmarshal hello data failed:", err)
		_ = conn.Close()
		return
	}

	msgNodeID := msg.NodeID

	nodeConnMapMu.Lock()
	nodeConnMap[msgNodeID] = conn
	nodeConnMapMu.Unlock()

	if hello.StableID != "" {
		types.PeerIDToStableID[msgNodeID] = hello.StableID
		types.StableIDToPeerID[hello.StableID] = msgNodeID
		log.Printf("Registered stableID %s for peerID %s\n", hello.StableID, msgNodeID)
	}

	writer := bufio.NewWriter(conn)

	welcome := Message{
		Type:   Msgwelcome,
		NodeID: e.MyID,
		MessageData: WelcomeMessage{
			NodeID: msgNodeID,
		},
	}

	jsonMsg, err := json.Marshal(welcome)
	if err != nil {
		log.Println("json marshal error in handleNewNode:", err)
		_ = conn.Close()
		return
	}

	_, err = writer.WriteString(string(jsonMsg) + "\n")
	if err != nil {
		log.Println("write welcome failed:", err)
		_ = conn.Close()
		return
	}

	if err := writer.Flush(); err != nil {
		log.Println("flush welcome failed:", err)
		_ = conn.Close()
		return
	}

	handleRestoreCabOrders(ps, e, msgNodeID, hello.StableID)

	go readLoop(conn, incomingTCP)
}

func ConnectToPrimary(ps *types.PeerState, port string, e *types.Elevator, incomingTCP chan Message) {
	for {
		if ps.PrimaryIP == "" {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if ps.PrimaryConn != nil {
			_ = ps.PrimaryConn.Close()
			ps.PrimaryConn = nil
		}

		primaryAddr := ps.PrimaryIP + ":" + port

		conn, err := net.Dial("tcp", primaryAddr)
		if err != nil {
			log.Println("Error connecting to primary:", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		log.Println("Connected to primary")

		ps.PrimaryConn = conn
		writer := bufio.NewWriter(conn)

		hello := Message{
			Type:   Msghello,
			NodeID: e.MyID,
			MessageData: HelloMessage{
				Role:     types.RoleToString(ps.Role),
				StableID: e.StableID,
			},
		}

		jsonMsg, err := json.Marshal(hello)
		if err != nil {
			log.Println("json marshal error in ConnectToPrimary:", err)
			_ = conn.Close()
			ps.PrimaryConn = nil
			return
		}

		_, err = writer.WriteString(string(jsonMsg) + "\n")
		if err != nil {
			log.Println("write hello failed:", err)
			_ = conn.Close()
			ps.PrimaryConn = nil
			time.Sleep(500 * time.Millisecond)
			continue
		}

		err = writer.Flush()
		if err != nil {
			log.Println("flush hello failed:", err)
			_ = conn.Close()
			ps.PrimaryConn = nil
			time.Sleep(500 * time.Millisecond)
			continue
		}

		go readLoop(conn, incomingTCP)
		break
	}
}
func SendTCP(receiverID string, message Message, ps *types.PeerState) {
	var connNode net.Conn

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Println("error in json conversion:", err)
		return
	}

	if ps.Role == types.RolePrimary {
		nodeConnMapMu.RLock()
		conn, exists := nodeConnMap[receiverID]
		nodeConnMapMu.RUnlock()

		if !exists {
			log.Println("No node in nodeConnMap:", receiverID)
			return
		}
		if conn == nil {
			log.Println("Connection is nil for receiver:", receiverID)
			return
		}
		connNode = conn
	} else {
		if ps.PrimaryConn == nil {
			log.Println("PrimaryConn is nil, cannot send")
			return
		}
		connNode = ps.PrimaryConn
	}

	writer := bufio.NewWriter(connNode)
	_, err = writer.WriteString(string(jsonMessage) + "\n")
	if err != nil {
		log.Println("WriteString failed:", err)
		return
	}
	if err := writer.Flush(); err != nil {
		log.Println("Flush failed:", err)
		return
	}
}

func HeartbeatTick(e *types.Elevator, ps *types.PeerState, d time.Duration, TCPHeartBeat chan<- Message) {
	tic := time.NewTicker(d)
	defer tic.Stop()

	for range tic.C {
		if ps.Role == types.RolePrimary {
			continue
		}

		heartbeat := HeartbeatMessage{
			CurrentFloor: e.CurrentFloor,
			State:        e.State,
			Dir:          e.Dir,
			CabRequests:  e.CabOrderMatrix[:],
			StableID:     e.StableID,
		}

		msg := Message{
			Type:        MsgHeartbeat,
			NodeID:      e.MyID,
			MessageData: heartbeat,
		}

		TCPHeartBeat <- msg
	}
}

func StartHeartbeatSender(ps *types.PeerState, heartbeatCh <-chan Message) {
	go func() {
		for msg := range heartbeatCh {
			if ps.PrimaryID != "" {
				SendTCP(ps.PrimaryID, msg, ps)
			}
		}
	}()
}
func handleRestoreCabOrders(ps *types.PeerState, e *types.Elevator, peerID string, stableID string) {
	if ps.Role != types.RolePrimary {
		return
	}

	if stableID == "" {
		return
	}

	cabs, ok := types.LostCabOrders[stableID]
	if !ok {
		return
	}

	log.Printf("Restoring cab orders for stableID %s to peerID %s: %v\n", stableID, peerID, cabs)

	messageData := RestoreCabOrdersMessage{
		NodeID: peerID,
		Cabs:   cabs,
	}

	buttonMessage := Message{
		Type:        MsgRestoreCabOrders,
		NodeID:      e.MyID,
		MessageData: messageData,
	}

	SendTCP(peerID, buttonMessage, ps)
	delete(types.LostCabOrders, stableID)
}
