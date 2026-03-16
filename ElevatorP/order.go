package ElevatorP

import (
	"heisprosjekt75/Driver-go/elevio"
	messagecomplete "heisprosjekt75/Messages/MessageComplete"
	"heisprosjekt75/types"
	"log"
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

func cabOrdersBelow(e *types.Elevator) bool {
	for f := e.CurrentFloor - 1; f >= 0; f-- {
		if e.CabOrderMatrix[f] {
			return true
		}
	}
	return false
}
func cabOrdersAbove(e *types.Elevator) bool {
	for f := e.CurrentFloor + 1; f < types.NumFloors; f++ {
		if e.CabOrderMatrix[f] {
			return true
		}
	}
	return false
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
			if cabOrdersBelow(e) && e.OrderDir == elevio.MD_Down {
				return elevio.MD_Down, types.Moving
			}
			return elevio.MD_Up, types.Moving
		}

		if orderBelow(e) {
			if cabOrdersAbove(e) && e.OrderDir == elevio.MD_Up {
				return elevio.MD_Up, types.Moving
			}
			return elevio.MD_Down, types.Moving
		}

		return elevio.MD_Stop, types.Idle

	case elevio.MD_Down:
		if orderBelow(e) {
			if cabOrdersAbove(e) && e.OrderDir == elevio.MD_Up {
				return elevio.MD_Up, types.Moving
			}
			return elevio.MD_Down, types.Moving
		}

		if orderAbove(e) {
			if cabOrdersBelow(e) && e.OrderDir == elevio.MD_Down {
				return elevio.MD_Down, types.Moving
			}
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
		return cabOrdersHere(e) || hallOrderUpHere(e) || (!orderAbove(e) && hallOrderDownHere(e))

	case elevio.MD_Down:
		return cabOrdersHere(e) || hallOrderDownHere(e) || (!orderBelow(e) && hallOrderUpHere(e))

	case elevio.MD_Stop:
		return cabOrdersHere(e) || hallOrderUpHere(e) || hallOrderDownHere(e)

	default:
		return false
	}
}

func setOrderDirAtStop(e *types.Elevator) {
	switch e.Dir {

	case elevio.MD_Up:
		if hallOrderUpHere(e) {
			e.OrderDir = elevio.MD_Up
		} else if cabOrdersHere(e) {
			e.OrderDir = elevio.MD_Up
		} else if !orderAbove(e) && hallOrderDownHere(e) {
			e.OrderDir = elevio.MD_Down
		}

	case elevio.MD_Down:
		if hallOrderDownHere(e) {
			e.OrderDir = elevio.MD_Down
		} else if cabOrdersHere(e) {
			e.OrderDir = elevio.MD_Down
		} else if !orderBelow(e) && hallOrderUpHere(e) {
			e.OrderDir = elevio.MD_Up
		}

	case elevio.MD_Stop:
		if hallOrderUpHere(e) {
			e.OrderDir = elevio.MD_Up
		} else if hallOrderDownHere(e) {
			e.OrderDir = elevio.MD_Down
		} else {
			e.OrderDir = elevio.MD_Stop
		}
	}
}

func shouldClearAtFloorImmediately(e *types.Elevator, btnFloor int, btnType elevio.ButtonType) bool {
	if elevio.GetFloor() == -1 {
		return false
	}

	return elevio.GetFloor() == btnFloor &&
		((e.Dir == elevio.MD_Up && btnType == elevio.BT_HallUp) ||
			(e.Dir == elevio.MD_Down && btnType == elevio.BT_HallDown) ||
			(e.Dir == elevio.MD_Stop) ||
			(btnType == elevio.BT_Cab))
}

/*
func clearAtCurrentFloor(e *types.Elevator, prevDir elevio.MotorDirection, ps *types.PeerState) {
	e.CabOrderMatrix[e.CurrentFloor] = false
	TurnOffCabLight(e.CurrentFloor)

	var btnTypeCleared []elevio.ButtonEvent

	if e.Mode == types.PrimaryBackup {
		if e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallUp] {
			e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallUp] = false
			btnTypeCleared = append(btnTypeCleared, elevio.ButtonEvent{
				Floor:  e.CurrentFloor,
				Button: elevio.BT_HallUp,
			})
		}

		if e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallDown] {
			e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallDown] = false
			btnTypeCleared = append(btnTypeCleared, elevio.ButtonEvent{
				Floor:  e.CurrentFloor,
				Button: elevio.BT_HallDown,
			})
		}

		for _, btn := range btnTypeCleared {
			messagecomplete.OrderCompleted(btn, e, ps)
		}
		return // ENDRET
	}

	switch prevDir {
	case elevio.MD_Up:
		if e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallUp] {
			e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallUp] = false
			TurnOffHallLight(elevio.BT_HallUp, e.CurrentFloor)
			btnTypeCleared = append(btnTypeCleared, elevio.ButtonEvent{
				Floor:  e.CurrentFloor,
				Button: elevio.BT_HallUp,
			})
		}

		if !orderAbove(e) && e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallDown] {
			e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallDown] = false
			TurnOffHallLight(elevio.BT_HallDown, e.CurrentFloor)
			btnTypeCleared = append(btnTypeCleared, elevio.ButtonEvent{
				Floor:  e.CurrentFloor,
				Button: elevio.BT_HallDown,
			})
		}

	case elevio.MD_Down:
		if e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallDown] {
			e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallDown] = false
			TurnOffHallLight(elevio.BT_HallDown, e.CurrentFloor)
			btnTypeCleared = append(btnTypeCleared, elevio.ButtonEvent{
				Floor:  e.CurrentFloor,
				Button: elevio.BT_HallDown,
			})
		}

		if !orderBelow(e) && e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallUp] {
			e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallUp] = false
			TurnOffHallLight(elevio.BT_HallUp, e.CurrentFloor)
			btnTypeCleared = append(btnTypeCleared, elevio.ButtonEvent{
				Floor:  e.CurrentFloor,
				Button: elevio.BT_HallUp,
			})
		}

	case elevio.MD_Stop:
		if e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallDown] {
			e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallDown] = false
			TurnOffHallLight(elevio.BT_HallDown, e.CurrentFloor)
			btnTypeCleared = append(btnTypeCleared, elevio.ButtonEvent{
				Floor:  e.CurrentFloor,
				Button: elevio.BT_HallDown,
			})
		}

		if e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallUp] { // ENDRET: ikke else if
			e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallUp] = false
			TurnOffHallLight(elevio.BT_HallUp, e.CurrentFloor)
			btnTypeCleared = append(btnTypeCleared, elevio.ButtonEvent{
				Floor:  e.CurrentFloor,
				Button: elevio.BT_HallUp,
			})
		}
	}

	if e.Mode == types.SingleElevator {
		for _, btn := range btnTypeCleared {
			_ = btn
		}
	}
}*/

func clearAtCurrentFloor(e *types.Elevator, prevDir elevio.MotorDirection, ps *types.PeerState) {

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
			messagecomplete.OrderCompleted(btn, e, ps)
		}
	}
}

func HandleAsignedOrder(e *types.Elevator, btnFloor int, btnButton elevio.ButtonType, doorStartTimerCh chan int, ps *types.PeerState) {
	if !e.HallOrderMatrix[btnFloor][btnButton] {
		log.Printf("Assigned order -> role:%v floor:%d button:%d\n", ps.Role, btnFloor, btnButton)
		AddOrder(e, btnFloor, btnButton)
	}

	if shouldClearAtFloorImmediately(e, btnFloor, btnButton) {
		switch btnButton {
		case elevio.BT_HallUp:
			e.OrderDir = elevio.MD_Up
		case elevio.BT_HallDown:
			e.OrderDir = elevio.MD_Down
		case elevio.BT_Cab:
			e.OrderDir = e.Dir
		}

		onDoorOpen(doorStartTimerCh, e, ps)
	} else {
		StartAction(e, doorStartTimerCh, ps)
	}
}

func SingleElevatorOrderRedelegation(e *types.Elevator, doorStartTimerCh chan int) {
	log.Println("Vi skal nå sjekke gjennom ordrelista")
	for f := 0; f < types.NumFloors; f++ {
		for b := 0; b < types.NumHallButtons; b++ {
			if types.FullOrderMatrix[f][b] {
				AddOrder(e, f, elevio.ButtonType(b))
				StartAction(e, doorStartTimerCh, &e.Ps)
				log.Println("Det finnes ordre på etg: ", f)
			} else {
				log.Println("Det finnes ingen ordre i etg: ", f)
			}
		}
	}
}
