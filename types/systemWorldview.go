package types

import (
	"heisprosjekt75/Driver-go/elevio"
	"time"
)

var FullOrderMatrix [NumFloors][NumHallButtons]bool
var HallOrderTimes [NumFloors][NumHallButtons]time.Time
var WorldView = make(map[string]ElevatorStatus)
var LostCabOrders = make(map[string][NumFloors]bool)
var PeerIDToStableID = make(map[string]string)
var StableIDToPeerID = make(map[string]string)
var CurrentAssignment = make(map[string]HallAssignment)

type HallAssignment [NumFloors][NumHallButtons]bool

type OrderTimestamp struct {
	CreatedAt time.Time
}

type ElevatorStatus struct {
	Floor       int                   `json:"floor"`
	Direction   elevio.MotorDirection `json:"direction"`
	State       ElevatorState         `json:"state"`
	CabRequests []bool                `json:"cabRequests"`
}

type ElevatorStateJSON struct {
	Behaviour   string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [][2]bool                    `json:"hallRequests"`
	States       map[string]ElevatorStateJSON `json:"states"`
}
