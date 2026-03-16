package tcp

import (
	"bufio"
	"encoding/json"
	"heisprosjekt75/types"
	"io"
	"log"
	"net"
	"strings"
	"time"
)


// Message struct og ulike meldinger må være definert et annet sted
// type Message {...}
// type HeartbeatMessage {...}
// type HelloMessage {...}
// type WelcomeMessage {...}
// const Msghello, Msgwelcome, MsgHeartbeat, etc. {...}

func readLoop(conn net.Conn, incomingTCP chan Message) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("Connection closed by client.")
			} else {
				log.Println("Message reading error in read loop:", err)
			}
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var msg Message
		err = json.Unmarshal([]byte(line), &msg)
		if err != nil {
			log.Println("json decode error in read loop:", err)
			continue
		}

		incomingTCP <- msg
	}
}

func StartPrimaryTCP(ps *types.PeerState, port string, incomingTCP chan Message, e *types.Elevator) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println("Error creating listener:", err)
		return
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Error accepting TCP connection:", err)
				continue
			}
			go handleNewNode(conn, incomingTCP, e, ps)
		}
	}()
}

func handleNewNode(conn net.Conn, incomingTCP chan Message, e *types.Elevator, ps *types.PeerState) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Message reading error in handleNewNode:", err)
		return
	}

	line = strings.TrimSpace(line)
	var msg Message
	if err := json.Unmarshal([]byte(line), &msg); err != nil {
		log.Println("json decode error in handleNewNode:", err)
		return
	}

	if msg.Type != Msghello {
		log.Println("expected Msghello, got", msg.Type)
		return
	}

	var hello HelloMessage
	helloBytes, _ := json.Marshal(msg.MessageData)
	if err := json.Unmarshal(helloBytes, &hello); err != nil {
		log.Println("unmarshal hello data failed:", err)
		return
	}

	nodeConnMapMu.Lock()
	nodeConnMap[msg.NodeID] = conn
	nodeConnMapMu.Unlock()

	if hello.StableID != "" {
		types.PeerIDToStableID[msg.NodeID] = hello.StableID
		types.StableIDToPeerID[hello.StableID] = msg.NodeID
		log.Printf("Registered stableID %s for peerID %s\n", hello.StableID, msg.NodeID)
	}

	writer := bufio.NewWriter(conn)

	welcome := Message{
		Type:   Msgwelcome,
		NodeID: e.MyID,
		MessageData: WelcomeMessage{
			NodeID: msg.NodeID,
		},
	}

	jsonMsg, _ := json.Marshal(welcome)
	writer.WriteString(string(jsonMsg) + "\n")
	writer.Flush()

	handleRestoreCabOrders(ps, e, msg.NodeID, hello.StableID)
	go readLoop(conn, incomingTCP)
}

func ConnectToPrimary(ps *types.PeerState, port string, e *types.Elevator, incomingTCP chan Message) {
	for {
		if ps.PrimaryIP == "" {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if ps.PrimaryConn != nil {
			ps.PrimaryConn.Close()
			ps.PrimaryConn = nil
		}

		conn, err := net.Dial("tcp", ps.PrimaryIP+":"+port)
		if err != nil {
			log.Println("Error connecting to primary:", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}

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

		jsonMsg, _ := json.Marshal(hello)
		writer.WriteString(string(jsonMsg) + "\n")
		writer.Flush()

		go readLoop(conn, incomingTCP)
		break
	}
}

func SendTCP(receiverID string, message Message, ps *types.PeerState) {
	var connNode net.Conn

	jsonMessage, _ := json.Marshal(message)

	if ps.Role == types.RolePrimary {
		nodeConnMapMu.RLock()
		conn, exists := nodeConnMap[receiverID]
		nodeConnMapMu.RUnlock()

		if !exists || conn == nil {
			log.Println("Cannot send: connection not found for", receiverID)
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
	writer.WriteString(string(jsonMessage) + "\n")
	writer.Flush()
}

func HeartbeatTick(e *types.Elevator, ps *types.PeerState, d time.Duration, TCPHeartBeat chan<- Message) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for range ticker.C {
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
	if ps.Role != types.RolePrimary || stableID == "" {
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
