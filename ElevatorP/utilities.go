package ElevatorP

import (
	"heisprosjekt75/Driver-go/elevio"
	messagecomplete "heisprosjekt75/Messages/MessageComplete"
	"heisprosjekt75/types"
)

func IsSingleElevatorMode(e *types.Elevator) bool {
	return e.Mode == types.SingleElevator
}

func shouldClearOppositeAtCurrentFloor(e *types.Elevator) bool {
	switch e.OrderDir {
	case elevio.MD_Up:
		return hallOrderDownHere(e) && !orderAbove(e)

	case elevio.MD_Down:
		return hallOrderUpHere(e) && !orderBelow(e)

	default:
		return false
	}

}

func clearOppositeAtAtCurrentFloor(e *types.Elevator, ps *types.PeerState) {
	var btn elevio.ButtonType
	var exist bool

	switch e.OrderDir {
	case elevio.MD_Up:
		btn = elevio.BT_HallDown
		exist = e.HallOrderMatrix[e.CurrentFloor][btn]

	case elevio.MD_Down:
		btn = elevio.BT_HallUp
		exist = e.HallOrderMatrix[e.CurrentFloor][btn]

	default:
		return
	}

	if !exist {
		return
	}

	e.HallOrderMatrix[e.CurrentFloor][btn] = false
	if IsSingleElevatorMode(e) {
		TurnOffHallLight(btn, e.CurrentFloor)
	}

	if !IsSingleElevatorMode(e) {
		btnEvent := elevio.ButtonEvent{Floor: e.CurrentFloor, Button: btn}
		messagecomplete.OrderCompleted(btnEvent, e, ps)
	}
}
