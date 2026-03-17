package messagestypes

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/types"
)

type MsgType int

const (
	MsgWorldView MsgType = iota
	MsgSnapshot
	MsgHallOrder
	MsgCompletedOrder
	MsgBackupHallOrderACK
	Msghello
	Msgwelcome
	MsgSetHallLights
	MsgTurnOffHallLights
	MsgStateSnapshot
	MsgRestoreCabOrders
)

type Message struct {
	Type        MsgType     `json:"type"`
	NodeID      string      `json:"nodeId"`
	MessageData interface{} `json:"messageData"`
}

type WorldViewMessage struct {
	CurrentFloor int                   `json:"currentFloor"`
	State        types.ElevatorState   `json:"state"`
	Dir          elevio.MotorDirection `json:"direction"`
	CabRequests  []bool                `json:"cabRequests"`
	StableID     string                `json:"stableId"`
}

type SnapshotHallOrdersMessage struct {
	Hall [types.NumFloors][types.NumHallButtons]bool `json:"hall"`
}

type StateSnapshotMessage struct {
	Hall             [types.NumFloors][types.NumHallButtons]bool
	WorldView        map[string]types.ElevatorStatus
	LostCabOrders    map[string][types.NumFloors]bool
	PeerIDToStableID map[string]string
	StableIDToPeerID map[string]string
	CabOrders        map[string][]bool
}

type RestoreCabOrdersMessage struct {
	NodeID string                `json:"nodeId"`
	Cabs   [types.NumFloors]bool `json:"cabs"`
}

type BackupHallOrderACK struct {
	Ack bool `json:"ack"`
}

type HelloMessage struct {
	Role     string
	StableID string
}

type WelcomeMessage struct {
	NodeID string `json:"nodeId"`
}

type HallOrderMessage struct {
	Floor  int               `json:"floor"`
	Button elevio.ButtonType `json:"button"`
}

type CompletedOrderMessage struct {
	Floor  int               `json:"floor"`
	Button elevio.ButtonType `json:"button"`
}

type HallLightsOnMessage struct {
	Floor  int               `json:"floor"`
	Button elevio.ButtonType `json:"button"`
}

type HallLightsOffMessage struct {
	Floor  int               `json:"floor"`
	Button elevio.ButtonType `json:"button"`
}
