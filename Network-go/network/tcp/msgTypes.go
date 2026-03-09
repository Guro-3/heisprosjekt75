package tcp

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/types"
	"net"
)

// leste meg opp på at når en skal sende struct over tcp så må en gjøre det om til jason, har prøvd meg frem er litt usikker
type MsgType int

const (
	MsgHeartbeat MsgType = iota
	MsgSnapshot
	MsgHallOrder
	MsgCompletedOrder
	MsgBackupHallOrderACK
	Msghello
	Msgwelcome
	
	
)


type Message struct {
	Type   MsgType         `json:"type"`
	NodeID string          `json:"nodeId"`
	MessageData interface{}     //`json:"data"`
}



type HeartbeatMessage struct {
	CurrentFloor int `json:"currentFloor"` 
	State types.ElevatorState `json:"state"` 
	Dir elevio.MotorDirection `json:"direction"` 
	CabRequests []bool `json:"cabRequests"` 
}

type SnapshotHallOrdersMessage struct {
	
	Hall [types.NumFloors][types.NumHallButtons]bool `json:"hall"`
}

type BackupHallOrderACK struct {
	Ack bool `json:"ack"`
	
}

type HelloMessage struct {
	Role string `json:"role"`
}


type WelcomeMessage struct {
	NodeID string `json:"role"`
}


type HallOrderMessage struct {
	Floor  int `json:"floor"`
	Button elevio.ButtonType `json:"button"` 
}

type CompletedOrderMessage struct {
	Floor  int `json:"floor"`
	Button elevio.ButtonType `json:"button"`
}

var nodeConnMap = make(map[string]net.Conn)




// per nå sender vi bare rene strings over nettet, fant en link kan se på det : https://agirlamonggeeks.com/convert-struct-to-json-string/
// vi må gjøre structene vår om til json og finne hvordan lese dem.

// når det er oppe å går, vi kan teste at knappe trykk som blir trykket kommer over nettet er neste steg å bruke utdelt cost funksjon slik at master kan best mulig delegere roller. og da er det mye testing som gjenstår tror jeg teste i forholde til spec