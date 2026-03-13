package ElevatorP

import (
	"heisprosjekt75/Driver-go/elevio"
	messagecomplete "heisprosjekt75/Messages/MessageComplete"
	"heisprosjekt75/types"
)

func AddOrder(e *types.Elevator, btnFloor int, btn elevio.ButtonType) {
	switch btn {
	case elevio.BT_Cab:
		e.CabOrderMatrix[btnFloor] = true
		SetCabLight(btnFloor)
	case elevio.BT_HallUp:
		e.HallOrderMatrix[btnFloor][elevio.BT_HallUp] = true
		if e.Mode == types.SingleElevator {
			SetHallLight(elevio.BT_HallUp, btnFloor)
		}
	case elevio.BT_HallDown:
		e.HallOrderMatrix[btnFloor][elevio.BT_HallDown] = true
		if e.Mode == types.SingleElevator {
			SetHallLight(elevio.BT_HallDown, btnFloor)
		}
	}
}

func cabOrdersHere(e *types.Elevator) bool {
	return e.CabOrderMatrix[e.CurrentFloor]
}

func hallOrderUpHere(e *types.Elevator) bool {
	return e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallUp]
}

func hallOrderDownHere(e *types.Elevator) bool {
	return e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallDown]
}

func orderBelow(e *types.Elevator) bool {
	for f := e.CurrentFloor - 1; f >= 0; f-- {
		for b := 0; b < types.NumHallButtons; b++ {
			if e.HallOrderMatrix[f][b] {
				return true
			}
		}

		if e.CabOrderMatrix[f] {
			return true
		}
	}
	return false
}

func orderAbove(e *types.Elevator) bool {
	for f := e.CurrentFloor + 1; f < types.NumFloors; f++ {
		for b := 0; b < types.NumHallButtons; b++ {
			if e.HallOrderMatrix[f][b] {
				return true
			}
		}
		if e.CabOrderMatrix[f] {
			return true
		}
	}
	return false
}

func chooseDirection(e *types.Elevator) (elevio.MotorDirection, types.ElevatorState) {
	switch e.Dir {

	case elevio.MD_Up:
		if orderAbove(e) {
			return elevio.MD_Up, types.Moving
		}

		if orderBelow(e) {
			return elevio.MD_Down, types.Moving
		}
		return elevio.MD_Stop, types.Idle

	case elevio.MD_Down:

		if orderBelow(e) {
			return elevio.MD_Down, types.Moving
		}
		if orderAbove(e) {
			return elevio.MD_Up, types.Moving
		}
		return elevio.MD_Stop, types.Idle

	case elevio.MD_Stop:

		if orderAbove(e) {
			return elevio.MD_Up, types.Moving
		}
		if orderBelow(e) {
			return elevio.MD_Down, types.Moving
		}
		return elevio.MD_Stop, types.Idle

	default:
		return elevio.MD_Stop, types.Idle
	}
}

func shouldStop(e *types.Elevator) bool {
	switch e.Dir {
	case elevio.MD_Up:
		return cabOrdersHere(e) || hallOrderUpHere(e) || !orderAbove(e)
	case elevio.MD_Down:
		return cabOrdersHere(e) || hallOrderDownHere(e) || !orderBelow(e)

	case elevio.MD_Stop:
		return true

	default:
		return true
	}
}

func shouldClearAtFloorImmediately(e *types.Elevator, btnFloor int, btnType elevio.ButtonType) bool {
	return e.CurrentFloor == btnFloor &&
		((e.Dir == elevio.MD_Up && btnType == elevio.BT_HallUp) ||
			(e.Dir == elevio.MD_Down && btnType == elevio.BT_HallDown) ||
			(e.Dir == elevio.MD_Stop) ||
			(btnType == elevio.BT_Cab))
}

func clearAtCurrentFloor(e *types.Elevator, prevDir elevio.MotorDirection, ps *types.PeerState) {

	e.CabOrderMatrix[e.CurrentFloor] = false
	TurnOffCabLight(e.CurrentFloor)

	var btnType elevio.ButtonType

	switch prevDir {
	case elevio.MD_Up:
		e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallUp] = false
		if e.Mode == types.SingleElevator {
			TurnOffHallLight(elevio.BT_HallUp, e.CurrentFloor)
		}
		btnType = elevio.BT_HallUp

		if !orderAbove(e) {
			e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallDown] = false
			if e.Mode == types.SingleElevator {
				TurnOffHallLight(elevio.BT_HallDown, e.CurrentFloor)
			}
			btnType = elevio.BT_HallDown
		}
	case elevio.MD_Down:

		e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallDown] = false
		if e.Mode == types.SingleElevator {
			TurnOffHallLight(elevio.BT_HallDown, e.CurrentFloor)
		}
		btnType = elevio.BT_HallDown

		if !orderBelow(e) {
			e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallUp] = false
			if e.Mode == types.SingleElevator {
				TurnOffHallLight(elevio.BT_HallUp, e.CurrentFloor)
			}
			btnType = elevio.BT_HallUp
		}
	}
	btn := elevio.ButtonEvent{
		Floor:  e.CurrentFloor,
		Button: btnType,
	}

	if e.Mode == types.PrimaryBackup {
		messagecomplete.OrderCompleted(btn, e, ps)
	}

}

func HandleAsignedOrder(e *types.Elevator, btnFloor int, btnButton elevio.ButtonType, doorStartTimerCh chan int, ps *types.PeerState) {
	if shouldClearAtFloorImmediately(e, btnFloor, btnButton) {
		onDoorOpen(doorStartTimerCh, e, ps)
		btn := elevio.ButtonEvent{Floor: btnFloor, Button: btnButton}
		messagecomplete.OrderCompleted(btn, e, ps)

	} else {
		AddOrder(e, btnFloor, btnButton)
		StartAction(e)
	}
}
