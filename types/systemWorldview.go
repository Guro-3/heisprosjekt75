package types

import (
	"heisprosjekt75/Driver-go/elevio"
)

type ElevatorStatus struct {
	Floor       int                   `json:"floor"`
	Direction   elevio.MotorDirection `json:"direction"`
	State       ElevatorState         `json:"state"`
	CabRequests []bool                `json:"cabRequests"`
}

var FullOrderMatrix [NumFloors][NumHallButtons]bool

var WorldView = make(map[string]ElevatorStatus)
var LostCabOrders = make(map[string][NumFloors]bool)
var PeerIDToStableID = make(map[string]string)
var StableIDToPeerID = make(map[string]string)

type HallAssignment [NumFloors][NumHallButtons]bool

var CurrentAssignment = make(map[string]HallAssignment)

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
