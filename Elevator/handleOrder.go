package Elevator

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/Messages/MessageComplete"
	"heisprosjekt75/types"
	"log"
)

func AddOrder(e *types.Elevator, btnFloor int, btnType elevio.ButtonType) {
	switch btnType {
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



func shouldClearOrderAtFloorImmediately(e *types.Elevator, btnFloor int, btnType elevio.ButtonType) bool {
	if elevio.GetFloor() == -1 {
		return false
	}

	return elevio.GetFloor() == btnFloor &&
		((e.Dir == elevio.MD_Up && btnType == elevio.BT_HallUp) ||
			(e.Dir == elevio.MD_Down && btnType == elevio.BT_HallDown) ||
			(e.Dir == elevio.MD_Stop) ||
			(btnType == elevio.BT_Cab))
}




func clearOrderAtCurrentFloor(e *types.Elevator, prevDir elevio.MotorDirection) {
	e.CabOrderMatrix[e.CurrentFloor] = false
	TurnOffCabLight(e.CurrentFloor)

	var btnCleared []elevio.ButtonEvent

	switch prevDir {
	case elevio.MD_Up:
		if e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallUp] {
			e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallUp] = false

			if IsSingleElevatorMode(e) {
				TurnOffHallLight(elevio.BT_HallUp, e.CurrentFloor)
			}

			btnCleared = append(btnCleared, elevio.ButtonEvent{
				Floor:  e.CurrentFloor,
				Button: elevio.BT_HallUp,
			})
		}
	case elevio.MD_Down:
		if e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallDown] {
			e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallDown] = false

			if IsSingleElevatorMode(e) {
				TurnOffHallLight(elevio.BT_HallDown, e.CurrentFloor)
			}

			btnCleared = append(btnCleared, elevio.ButtonEvent{
				Floor:  e.CurrentFloor,
				Button: elevio.BT_HallDown,
			})
		}
	case elevio.MD_Stop:
		if e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallUp] {
			e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallUp] = false

			if IsSingleElevatorMode(e) {
				TurnOffHallLight(elevio.BT_HallUp, e.CurrentFloor)
			}

			btnCleared = append(btnCleared, elevio.ButtonEvent{
				Floor:  e.CurrentFloor,
				Button: elevio.BT_HallUp,
			})

		} else if e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallDown] {
			e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallDown] = false

			if IsSingleElevatorMode(e) {
				TurnOffHallLight(elevio.BT_HallDown, e.CurrentFloor)
			}

			btnCleared = append(btnCleared, elevio.ButtonEvent{
				Floor:  e.CurrentFloor,
				Button: elevio.BT_HallDown,
			})
		}
	}
	if !IsSingleElevatorMode(e) {
		for _, btn := range btnCleared {
			messagecomplete.OrderCompleted(btn, e)
		}
	}
}


func clearOppositeOrderAtCurrentFloor(e *types.Elevator) {
	var btnType elevio.ButtonType
	var exist bool

	switch e.OrderDir {
	case elevio.MD_Up:
		btnType = elevio.BT_HallDown
		exist = e.HallOrderMatrix[e.CurrentFloor][btnType]

	case elevio.MD_Down:
		btnType = elevio.BT_HallUp
		exist = e.HallOrderMatrix[e.CurrentFloor][btnType]

	default:
		return
	}

	if !exist {
		return
	}

	e.HallOrderMatrix[e.CurrentFloor][btnType] = false
	if IsSingleElevatorMode(e) {
		TurnOffHallLight(btnType, e.CurrentFloor)
	} else {
		btn := elevio.ButtonEvent{Floor: e.CurrentFloor, Button: btnType}
		messagecomplete.OrderCompleted(btn, e)
	}
}




func HandleAssignedOrder(e *types.Elevator, btnFloor int, btnType elevio.ButtonType, doorStartTimerCh chan int) {
	AddOrder(e, btnFloor, btnType)

	if shouldClearOrderAtFloorImmediately(e, btnFloor, btnType) {
		switch btnType {
		case elevio.BT_HallUp:
			e.OrderDir = elevio.MD_Up
		case elevio.BT_HallDown:
			e.OrderDir = elevio.MD_Down
		case elevio.BT_Cab:
			e.OrderDir = e.Dir
		}

		DoorOpen(doorStartTimerCh, e)
	} else {
		FsmStartAction(e, doorStartTimerCh)
	}
}




func SingleElevatorOrderRedelegation(e *types.Elevator, doorStartTimerCh chan int) {
	log.Println("Vi skal nå sjekke gjennom ordrelista")
	for f := 0; f < types.NumFloors; f++ {
		for b := 0; b < types.NumHallButtons; b++ {
			if types.FullOrderMatrix[f][b] {
				AddOrder(e, f, elevio.ButtonType(b))
				FsmStartAction(e, doorStartTimerCh)
				log.Println("Det finnes ordre på etg: ", f)
			} else {
				log.Println("Det finnes ingen ordre i etg: ", f)
			}
		}
	}
}
