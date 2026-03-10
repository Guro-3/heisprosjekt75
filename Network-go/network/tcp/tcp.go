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
			log.Println("Message reading error: %s\n", err)
			return
		}

		line = strings.TrimSpace(line)
		var msg Message
		err = json.Unmarshal([]byte(line), &msg)
		if err != nil {
			log.Println("json decode error: %v\n", err)
			continue
		}
		incomingTCP <- msg
	}
}

func StartPrimaryTCP(ps *types.PeerState, port string, incomingTCP chan Message, e *types.Elevator) {

	log.Println("starting primary")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println("Error with listning object")
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Error with listning object")
				continue
			}
			go handleNewNode(conn, incomingTCP, e)
		}
	}()
}

func handleNewNode(conn net.Conn, incomingTCP chan Message, e *types.Elevator) {
	log.Println("går inn i go handle")
	reader := bufio.NewReader(conn)

	line, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Message reading error: %v\n", err)
		_ = conn.Close()
		return
	}

	line = strings.TrimSpace(line)

	var msg Message
	err = json.Unmarshal([]byte(line), &msg)
	if err != nil {
		log.Printf("json decode error: %v\n", err)
		conn.Close()
		return
	}

	if msg.Type != Msghello {
		log.Printf("expected Msghello, got %v\n", msg.Type)
		conn.Close()
		return
	}

	// NB: nodeConnMap må eksistere som global map i pakken deres
	msgNodeID := msg.NodeID
	nodeConnMap[msgNodeID] = conn

	log.Printf("conn to node %v\n:", nodeConnMap[msgNodeID])
	log.Printf("Node connected to Primary %s\n:", msgNodeID)

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
		log.Println("json marshal error:", err)
		return
	}

	writer.WriteString(string(jsonMsg) + "\n")
	writer.Flush()

	go readLoop(conn, incomingTCP)
}

func ConnectToPrimary(ps *types.PeerState, port string, e *types.Elevator, incomingTCP chan Message) {

	for {
		log.Printf("DEBUG ConnectToPrimary: PrimaryID=%q PrimaryIP=%q port=%q", ps.PrimaryID, ps.PrimaryIP, port)

		if ps.PrimaryIP == "" {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		primaryAddr := ps.PrimaryIP + ":" + port

		conn, err := net.Dial("tcp", primaryAddr)
		if err != nil {
			log.Println("Error connecting to primary:", err)
			time.Sleep(5 * time.Second)
			continue
		}

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
			log.Println("json marshal error:", err)
			return
		}

		writer.WriteString(string(jsonMsg) + "\n")
		writer.Flush()

		log.Println("Connected to primary")

		go readLoop(conn, incomingTCP)
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
		log.Printf("Node connected to Primary recieverID: %s\n:", recieverID)
		log.Printf("conn to node %v\n:", nodeConnMap[recieverID])

		if !excist {
			log.Printf("No nodes in nodeConnMap %s\n", recieverID)
			return
		} else if conn == nil {
			log.Printf("No receiver in conn%s\n", recieverID)
			return
		}
		connNode = conn
	} else {
		if ps.PrimaryConn == nil {
			log.Printf("PrimaryConn is nil, cannot send\n")
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
