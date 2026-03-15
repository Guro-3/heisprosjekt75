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

func readLoop(conn net.Conn, incomingTCP chan Message, ps *types.PeerState, connectedNodeID string) {
	defer func() {
		_ = conn.Close()

		// Hvis dette er connectionen til primary, nullstill den når den dør
		if ps != nil && ps.PrimaryConn == conn {
			ps.PrimaryConn = nil
			log.Println("Primary connection lost -> ps.PrimaryConn = nil")
		}

		// Hvis dette er en node-connection hos primary, fjern den fra map
		if connectedNodeID != "" {
			if existingConn, exists := nodeConnMap[connectedNodeID]; exists && existingConn == conn {
				delete(nodeConnMap, connectedNodeID)
				log.Println("Removed node from nodeConnMap:", connectedNodeID)
			}
		}
	}()

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
			go handleNewNode(conn, incomingTCP, e)
		}
	}()
}

func handleNewNode(conn net.Conn, incomingTCP chan Message, e *types.Elevator) {
	reader := bufio.NewReader(conn)

	line, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Message reading error in handleNewNode: ", err)
		_ = conn.Close()
		return
	}

	line = strings.TrimSpace(line)

	var msg Message
	err = json.Unmarshal([]byte(line), &msg)
	if err != nil {
		log.Println("json decode error in handleNewNode: ", err)
		conn.Close()
		return
	}

	if msg.Type != Msghello {
		log.Println("expected Msghello, got ", msg.Type)
		conn.Close()
		return
	}

	// NB: nodeConnMap må eksistere som global map i pakken deres
	msgNodeID := msg.NodeID
	nodeConnMap[msgNodeID] = conn

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
		return
	}

	_, err = writer.WriteString(string(jsonMsg) + "\n")
	if err != nil {
		log.Println("write error in handleNewNode:", err)
		delete(nodeConnMap, msgNodeID)
		_ = conn.Close()
		return
	}

	err = writer.Flush()
	if err != nil {
		log.Println("flush error in handleNewNode:", err)
		delete(nodeConnMap, msgNodeID)
		_ = conn.Close()
		return
	}

	go readLoop(conn, incomingTCP, nil, msgNodeID)
}

func ConnectToPrimary(ps *types.PeerState, port string, e *types.Elevator, incomingTCP chan Message) {

	for {
		if ps.Role == types.RolePrimary {
			return
		}

		// Hvis vi allerede er koblet til, stopp
		if ps.PrimaryConn != nil {
			return
		}

		if ps.PrimaryIP == "" {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		primaryAddr := ps.PrimaryIP + ":" + port

		conn, err := net.Dial("tcp", primaryAddr)
		if err != nil {
			log.Println("Error connecting to primary:", err)
			time.Sleep(1 * time.Second)
			continue
		}
		log.Println("Connected to primary")

		ps.PrimaryConn = conn
		writer := bufio.NewWriter(conn)

		hello := Message{
			Type:   Msghello,
			NodeID: e.MyID,
			MessageData: HelloMessage{
				Role: types.RoleToString(ps.Role),
			},
		}

		jsonMsg, err := json.Marshal(hello)
		if err != nil {
			log.Println("json marshal error in ConnectToPrimary:", err)
			return
		}

		_, err = writer.WriteString(string(jsonMsg) + "\n")
		if err != nil {
			log.Println("write error:", err)
			_ = conn.Close()
			ps.PrimaryConn = nil
			time.Sleep(1 * time.Second)
			continue
		}

		err = writer.Flush()
		if err != nil {
			log.Println("flush error:", err)
			_ = conn.Close()
			ps.PrimaryConn = nil
			time.Sleep(1 * time.Second)
			continue
		}

		go readLoop(conn, incomingTCP, ps, "")
		break
	}
}

func SendTCP(recieverID string, message Message, ps *types.PeerState) {
	var connNode net.Conn

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Println("error in json convertion")
	}

	if ps.Role == types.RolePrimary {
		conn, excist := nodeConnMap[recieverID]

		if !excist {
			log.Println("No nodes in nodeConnMap ", recieverID)
			return
		} else if conn == nil {
			log.Println("No receiver in conn ", recieverID)
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
