package ElevatorP

import (
	"heisprosjekt75/Driver-go/elevio"	
)

func addOrder(e *Elevator, btnFloor int, btn elevio.ButtonType) {
	switch btn {
	case elevio.BT_Cab:
		e.CabOrderMatrix[btnFloor][0] = true
		SeCabLight(btnFloor)
	case elevio.BT_HallUp:
		e.HallorderMatrix[btnFloor][elevio.BT_HallUp] = true
		SetHallLight(elevio.BT_HallUp, btnFloor)
	case elevio.BT_HallDown:
		e.HallorderMatrix[btnFloor][elevio.BT_HallDown] = true
		SetHallLight(elevio.BT_HallDown, btnFloor)
	}
}


func cabOrdersHere(e *Elevator) bool {
	return e.CabOrderMatrix[e.CurrentFloor][0]
}

func hallOrderUpHere(e *Elevator) bool {
	return e.HallorderMatrix[e.CurrentFloor][elevio.BT_HallUp]
}

func hallOrderDownHere(e *Elevator) bool {
	return e.HallorderMatrix[e.CurrentFloor][elevio.BT_HallDown]
}


func orderBelow(e *Elevator) bool {
	for f := e.CurrentFloor- 1; f >= 0; f-- {
		for b := 0; b < numHallButtons; b++ {
			if e.HallorderMatrix[f][b] {
				return true
			}
		}
		
		if e.CabOrderMatrix[f][0] {
			return true
		}
	}
	return false
}


func orderAbove(e *Elevator) bool {
	for f := e.CurrentFloor+ 1; f < NUMFloors; f++ {
		for b := 0; b < numHallButtons; b++ {
			if e.HallorderMatrix[f][b] {
				return true
			}
		}
		if e.CabOrderMatrix[f][0] {
			return true
		}
	}
	return false
}


func chooseDirection(e *Elevator) (elevio.MotorDirection, elevatorState) {
	switch e.Dir{

	case elevio.MD_Up:
		if orderAbove(e) {
			return elevio.MD_Up, Moving
		}

		if orderBelow(e) {
			return elevio.MD_Down, Moving
		}
		return elevio.MD_Stop, Idle

	case elevio.MD_Down:
	
		if orderBelow(e) {
			return elevio.MD_Down, Moving
		}
		if orderAbove(e) {
			return elevio.MD_Up, Moving
		}
		return elevio.MD_Stop, Idle

	case elevio.MD_Stop:
		
		if orderAbove(e) {
			return elevio.MD_Up, Moving
		}
		if orderBelow(e) {
			return elevio.MD_Down, Moving
		}
		return elevio.MD_Stop, Idle

	default:
		return elevio.MD_Stop, Idle
	}
}


func shouldStop(e *Elevator) bool {
	switch e.Dir{
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


func shouldClearAtFloorImmediately(e *Elevator, btnFloor int, btnType elevio.ButtonType) bool {
	return e.CurrentFloor == btnFloor &&
		((e.Dir == elevio.MD_Up && btnType == elevio.BT_HallUp) ||
			(e.Dir == elevio.MD_Down && btnType == elevio.BT_HallDown) ||
			(e.Dir == elevio.MD_Stop) ||
			(btnType == elevio.BT_Cab))
}


func clearAtCurrentFloor(e *Elevator, prevDir elevio.MotorDirection) {

	e.CabOrderMatrix[e.CurrentFloor][0] = false
	TurnOffCabLight(e.CurrentFloor)

	switch prevDir{
	case elevio.MD_Up:
		e.HallorderMatrix[e.CurrentFloor][elevio.BT_HallUp] = false
		TurnOffHallLight(elevio.BT_HallUp, e.CurrentFloor)
		
		if !orderAbove(e) {
			e.HallorderMatrix[e.CurrentFloor][elevio.BT_HallDown] = false
			TurnOffHallLight(elevio.BT_HallDown, e.CurrentFloor)
		}
	case elevio.MD_Down:
		
		e.HallorderMatrix[e.CurrentFloor][elevio.BT_HallDown] = false
		TurnOffHallLight(elevio.BT_HallDown, e.CurrentFloor)
		
		if !orderBelow(e)  {
			e.HallorderMatrix[e.CurrentFloor][elevio.BT_HallUp] = false
			TurnOffHallLight(elevio.BT_HallUp, e.CurrentFloor)
		}
	}
}
