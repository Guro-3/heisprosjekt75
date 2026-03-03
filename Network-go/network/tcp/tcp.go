package tcp

import (
	"bufio"
	"heisprosjekt75/ElevatorP"
	"heisprosjekt75/RoleManager"
	"log"
	"net"
	"strings"
	"fmt"
)

func readLoop(conn net.Conn, incomingTCP chan string) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n') //hvor langt skal en string være? dette må skrives inn i parantesen
		if err != nil {
			log.Println("Message reading error: %s\n", err)
			return
		}
		incomingTCP <- strings.TrimSpace(message)
	}
}

func StartPrimaryTCP(ps *RoleManager.PeerState, port string, incomingTCP chan string) {
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

func handleNewNode(conn net.Conn, incomingTCP chan string) {

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

func ConnectToPrimary(ps *RoleManager.PeerState, port string, e *ElevatorP.Elevator, incomingTCP chan string) {
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
	msg := fmt.Sprintf("HELLO|%s|%s", e.MyID, RoleManager.RoleToString(ps.Role))
	writer.WriteString(msg + "\n")
	writer.Flush()
	log.Println("ID connected to primary: ", e.MyID)
	go readLoop(conn, incomingTCP)
}

func SendTCP(recieverID string, message string, ps *RoleManager.PeerState) {
	var connNode net.Conn

	if ps.Role == RoleManager.RolePrimary {
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
	writer.WriteString(message + "\n")
	writer.Flush()

}
