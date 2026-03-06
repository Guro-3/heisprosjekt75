package tcp

import (
	"bufio"
	"encoding/json"
	"fmt"
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

func StartPrimaryTCP(ps *types.PeerState, port string, incomingTCP chan Message) {
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
			go handleNewNode(conn, incomingTCP)
		}
	}()
}

func handleNewNode(conn net.Conn, incomingTCP chan Message) {

	reader := bufio.NewReader(conn)

	msgNodeID, err := reader.ReadString('\n') //hvor langt skal en string være? dette må skrives inn i parantesen
	if err != nil {
		log.Printf("Message reading error: %v\n", err)
		conn.Close()
		return
	}
	msgNodeID = strings.TrimSpace(msgNodeID)
	nodeConnMap[msgNodeID] = conn
	log.Printf("Node connected to Primary %s\n:", msgNodeID)
	writer := bufio.NewWriter(conn)
	writer.WriteString("WELCOME_" + msgNodeID + "\n")
	writer.Flush()

	go readLoop(conn, incomingTCP)
}

func ConnectToPrimary(ps *types.PeerState, port string, e *types.Elevator, incomingTCP chan Message) {
	primaryAddr := ps.PrimaryIP + ":" + port

	conn, err := net.Dial("tcp", primaryAddr)
	if err != nil {
		log.Println("Error accepting conn; ", err)
		return
	}

	writer := bufio.NewWriter(conn)
	ps.PrimaryConn = conn
	writer.WriteString(e.MyID + "\n")
	writer.Flush()
	msg := fmt.Sprintf("HELLO|%s|%s", e.MyID, types.RoleToString(ps.Role))
	writer.WriteString(msg + "\n")
	writer.Flush()
	log.Println("ID connected to primary: ", e.MyID)
	go readLoop(conn, incomingTCP)
}

func SendTCP(recieverID string, message Message, ps *types.PeerState) {
	var connNode net.Conn

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Println("error in json convertion")
	}

	if ps.Role == types.RolePrimary {
		conn, excist := nodeConnMap[recieverID]
		if !excist || conn == nil {
			log.Printf("No conn for receiver %s\n", recieverID)
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
	if ps.Role == types.RolePrimary {
		return
	}

	tic := time.NewTicker(d)
	defer tic.Stop()

	for range tic.C {
		cabRequests := types.MatrixToSlice(types.NumFloors, types.NumCabButtons, func(f, b int) bool { return e.CabOrderMatrix[f][b] })

		heartbeat := HeartbeatMessage{
			CurrentFloor: e.CurrentFloor,
			State:        e.State,
			Dir:          e.Dir,
			CabRequests:  cabRequests,
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
