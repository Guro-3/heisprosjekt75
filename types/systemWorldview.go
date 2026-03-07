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
var NodeOrderMap = make(map[string][]elevio.ButtonEvent)
var WorldView = make(map[string]ElevatorStatus)

type ElevatorStateJSON struct {
	Behaviour   string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

func WorldViewToJSON(world map[string]ElevatorStatus) map[string]ElevatorStateJSON {
	result := make(map[string]ElevatorStateJSON)
	for id, e := range world {
		dir := "stop"
		switch e.Direction {
		case elevio.MD_Up:
			dir = "up"
		case elevio.MD_Down:
			dir = "down"
		case elevio.MD_Stop:
			dir = "stop"
		}

		state := "idle"
		switch e.State {
		case Idle:
			state = "idle"
		case Moving:
			state = "moving"
		case DoorOpen:
			state = "doorOpen"
		}

		result[id] = ElevatorStateJSON{
			Behaviour:   state,
			Floor:       e.Floor,
			Direction:   dir,
			CabRequests: e.CabRequests,
		}
	}
	return result
}

type HRAInput struct {
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[string]ElevatorStateJSON `json:"states"`
}


func UpdateMyState(e *Elevator) {
    WorldView[e.MyID] = ElevatorStatus{
        Floor:       e.CurrentFloor,
        Direction:   e.Dir,
        State:       e.State,
        CabRequests: e.CabOrderMatrix[:],
    }
}
