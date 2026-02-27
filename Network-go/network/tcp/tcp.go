package tcp

import (
	"bufio"
	"heisprosjekt75/ElevatorP"
	"heisprosjekt75/RoleManager"
	"log"
	"net"
	"strings"
	"strconv"
)

func readLoop(conn net.Conn, incomingTCP chan string) {
	go func(){
		defer conn.Close()

		reader := bufio.NewReader(conn)
		for {
			message, msgErr := reader.ReadString('\n') //hvor langt skal en string være? dette må skrives inn i parantesen
			if msgErr != nil {
				log.Println("Message reading error: %s\n", msgErr)
				return
			}
			incomingTCP <- message
		}
	}()
}


func StartPrimaryTCP(ps *RoleManager.PeerState,port string, incomingTCP chan string) {
	primaryAddr := ps.PrimaryIP +":"+ port
	listener, err := net.Listen("tcp", primaryAddr)
	if err != nil{
		log.Println("Error with listning object")
	}
	go func(){
	for{
		conn, err := listener.Accept()
		if err != nil{
			log.Println("Error with listning object")
			continue
		}
		go handleNewNode(conn, incomingTCP)
	}
	}()
}

func handleNewNode(conn net.Conn, incomingTCP chan string) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	
	msgNodeID, msgErr := reader.ReadString('\n') //hvor langt skal en string være? dette må skrives inn i parantesen
	if msgErr != nil {
		log.Println("Message reading error: %s\n", msgErr)
		return
	}
	msgNodeID = strings.TrimSpace(msgNodeID)
	nodeConnMap[msgNodeID] = conn
	log.Println("Node connected to Primary %s\n:", msgNodeID)
	readLoop(conn, incomingTCP)
}

func ConnectToPrimary(ps *RoleManager.PeerState,port string, e *ElevatorP.Elevator, incomingTCP chan string) {
	intPort, _ := strconv.Atoi(port)
	primaryAddr := &net.TCPAddr{
		IP: net.IP(ps.PrimaryIP),
		Port: intPort,
	}
	conn, err := net.DialTCP("tcp", nil, primaryAddr)
	if err != nil {
		log.Println("Error accepting conn; ", err)
	}

	writer := bufio.NewWriter(conn)
	
	writer.WriteString(e.MyID + "\n")
	log.Println("ID connected to primary: ", e.MyID)
	readLoop(conn, incomingTCP)
}

func sendTCP(recieverID string, message string, ps RoleManager.PeerState) {
	if ps.Role == RoleManager.RolePrimary{
		connNode := nodeConnMap[recieverID] 
		writer := bufio.NewWriter(connNode)
		writer.WriteString(message + "\n")
	} else {
		connPrimary := ps.PrimaryConn
		writer := bufio.NewWriter(connPrimary)
		writer.WriteString(message + "\n")
	}
}

