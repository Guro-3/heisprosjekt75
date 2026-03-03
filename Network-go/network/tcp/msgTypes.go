package tcp

import (
	"encoding/json"
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/RoleManager"
	"net"
)

// leste meg opp på at når en skal sende struct over tcp så må en gjøre det om til jason, har prøvd meg frem er litt usikker
type MsgType int

const (
	MsgHeartbeat MsgType = 0
	MsgSnapshot
	MsgHallOrder
	MsgCompletedOrder
)


type Message struct {
	Type   MsgType         `json:"type"`
	NodeID string          `json:"nodeId"`
	Data   json.RawMessage `json:"data"`
}

// PAYLOADS:

type HeartbeatMessage struct {
	Role       RoleManager.ElevatorRole `json:"role"` 
	Floor      int                     `json:"floor"`
	CurrentDir elevio.MotorDirection   `json:"currentDir"`
}

type SnapshotHallOrdersMessage struct {
	
	Hall [][]bool `json:"hall"`
}


type HallOrderMessage struct {
	Floor  int `json:"floor"`
	Button int `json:"button"` 
}

type CompletedOrderMessage struct {
	Floor  int `json:"floor"`
	Button int `json:"button"`
}

var nodeConnMap = make(map[string]net.Conn)




// per nå sender vi bare rene strings over nettet, fant en link kan se på det : https://agirlamonggeeks.com/convert-struct-to-json-string/
// vi må gjøre structene vår om til json og finne hvordan lese dem.

// når det er oppe å går, vi kan teste at knappe trykk som blir trykket kommer over nettet er neste steg å bruke utdelt cost funksjon slik at master kan best mulig delegere roller. og da er det mye testing som gjenstår tror jeg teste i forholde til spec