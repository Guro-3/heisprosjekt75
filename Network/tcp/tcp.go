package tcp

import (
	"bufio"
	"encoding/json"
	"errors"
	messagestypes "heisprosjekt75/Messages/MessageTypes"
	"heisprosjekt75/types"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	nodeConnMap   = make(map[string]net.Conn)
	nodeConnMapMu sync.RWMutex
)

func tcpReadLoop(conn net.Conn, incomingTCP chan messagestypes.Message) {

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

		var msg messagestypes.Message
		err = json.Unmarshal([]byte(line), &msg)
		if err != nil {
			log.Println("json decode error in read loop:", err)
			continue
		}
		incomingTCP <- msg
	}
}

func TcpStartPrimary(port string, incomingTCP chan messagestypes.Message, e *types.Elevator) {

	if e.Ps.PrimaryListener != nil {
		return
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println("Error creating listener:", err)
		return
	}
	e.Ps.PrimaryListener = listener
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}
				log.Println("Error accepting TCP connection:", err)
				break
			}
			go tcpHandleNewNode(conn, incomingTCP, e)
		}
	}()
}

func tcpHandleNewNode(conn net.Conn, incomingTCP chan messagestypes.Message, e *types.Elevator) {
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Message reading error in handleNewNode:", err)
		return
	}

	line = strings.TrimSpace(line)
	var msg messagestypes.Message
	if err := json.Unmarshal([]byte(line), &msg); err != nil {
		log.Println("json decode error in handleNewNode:", err)
		return
	}

	if msg.Type != messagestypes.Msghello {
		log.Println("expected Msghello, got", msg.Type)
		return
	}

	var hello messagestypes.HelloMessage
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
	}

	writer := bufio.NewWriter(conn)
	welcome := messagestypes.Message{
		Type:        messagestypes.Msgwelcome,
		NodeID:      e.MyID,
		MessageData: messagestypes.WelcomeMessage{NodeID: msg.NodeID},
	}

	jsonMsg, _ := json.Marshal(welcome)
	writer.WriteString(string(jsonMsg) + "\n")
	writer.Flush()

	handleRestoreCabOrders(e, msg.NodeID, hello.StableID)
	go tcpReadLoop(conn, incomingTCP)
}

func TcpConnectToPrimary(port string, e *types.Elevator, incomingTCP chan messagestypes.Message) {
	for {
		if e.Ps.PrimaryIP == "" {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		/*if e.Ps.PrimaryConn != nil {
			e.Ps.PrimaryConn.Close()
			e.Ps.PrimaryConn = nil
		}*/

		conn, err := net.Dial("tcp", e.Ps.PrimaryIP+":"+port)
		if err != nil {
			log.Println("Error connecting to primary:", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		e.Ps.PrimaryConn = conn
		writer := bufio.NewWriter(conn)

		hello := messagestypes.Message{
			Type:        messagestypes.Msghello,
			NodeID:      e.MyID,
			MessageData: messagestypes.HelloMessage{Role: types.TypesRoleToString(e.Ps.Role), StableID: e.StableID},
		}

		jsonMsg, _ := json.Marshal(hello)
		writer.WriteString(string(jsonMsg) + "\n")
		writer.Flush()

		go tcpReadLoop(conn, incomingTCP)
		return
	}
}

func SendTCP(receiverID string, message messagestypes.Message, ps *types.PeerState) {
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
