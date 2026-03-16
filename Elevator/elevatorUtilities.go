package Elevator

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/types"
)

func IsSingleElevatorMode(e *types.Elevator) bool {
	return e.Mode == types.SingleElevator
}


func chooseDirection(e *types.Elevator) (elevio.MotorDirection, types.ElevatorState) {
	switch e.Dir {

	case elevio.MD_Up:
		if checkOrderAbove(e) {
			if checkCabOrdersBelow(e) && e.OrderDir == elevio.MD_Down {
				return elevio.MD_Down, types.Moving
			}
			return elevio.MD_Up, types.Moving
		}

		if checkOrderBelow(e) {
			if checkCabOrdersAbove(e) && e.OrderDir == elevio.MD_Up {
				return elevio.MD_Up, types.Moving
			}
			return elevio.MD_Down, types.Moving
		}

		return elevio.MD_Stop, types.Idle

	case elevio.MD_Down:
		if checkOrderBelow(e) {
			if checkCabOrdersAbove(e) && e.OrderDir == elevio.MD_Up {
				return elevio.MD_Up, types.Moving
			}
			return elevio.MD_Down, types.Moving
		}

		if checkOrderAbove(e) {
			if checkCabOrdersBelow(e) && e.OrderDir == elevio.MD_Down {
				return elevio.MD_Down, types.Moving
			}
			return elevio.MD_Up, types.Moving
		}

		return elevio.MD_Stop, types.Idle

	case elevio.MD_Stop:
		if checkOrderAbove(e) {
			return elevio.MD_Up, types.Moving
		}
		if checkOrderBelow(e) {
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
		return checkCabOrdersHere(e) || checkHallOrderUpHere(e) || (!checkOrderAbove(e) && checkHallOrderDownHere(e))

	case elevio.MD_Down:
		return checkCabOrdersHere(e) || checkHallOrderDownHere(e) || (!checkOrderBelow(e) && checkHallOrderUpHere(e))

	case elevio.MD_Stop:
		return checkCabOrdersHere(e) || checkHallOrderUpHere(e) || checkHallOrderDownHere(e)

	default:
		return false
	}
}


func shouldClearOppositeOrderAtCurrentFloor(e *types.Elevator) bool {
	switch e.OrderDir {
	case elevio.MD_Up:
		return checkHallOrderDownHere(e) && !checkOrderAbove(e)

	case elevio.MD_Down:
		return checkHallOrderUpHere(e) && !checkOrderBelow(e)

	default:
		return false
	}

}


func setOrderDirAtStop(e *types.Elevator) {
	switch e.Dir {

	case elevio.MD_Up:
		if checkHallOrderUpHere(e) {
			e.OrderDir = elevio.MD_Up
		} else if checkCabOrdersHere(e) {
			e.OrderDir = elevio.MD_Up
		} else if !checkOrderAbove(e) && checkHallOrderDownHere(e) {
			e.OrderDir = elevio.MD_Down
		}

	case elevio.MD_Down:
		if checkHallOrderDownHere(e) {
			e.OrderDir = elevio.MD_Down
		} else if checkCabOrdersHere(e) {
			e.OrderDir = elevio.MD_Down
		} else if !checkOrderBelow(e) && checkHallOrderUpHere(e) {
			e.OrderDir = elevio.MD_Up
		}

	case elevio.MD_Stop:
		if checkHallOrderUpHere(e) {
			e.OrderDir = elevio.MD_Up
		} else if checkHallOrderDownHere(e) {
			e.OrderDir = elevio.MD_Down
		} else {
			e.OrderDir = elevio.MD_Stop
		}
	}
}