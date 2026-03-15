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
		log.Println("Stop condition: ", e.StopCond)
		if orderAbove(e) {
			if cabOrdersBelow(e) && e.StopCond == types.DownOrder {
				return elevio.MD_Down, types.Moving
			} else {
				log.Println("chooseDir, MD_Up, OrderAbove()")
				return elevio.MD_Up, types.Moving
			}

		}

		if orderBelow(e) {
			if cabOrdersAbove(e) && e.StopCond == types.UpOrder {
				return elevio.MD_Up, types.Moving
			} else {
				log.Println("chooseDir, MD_Up, OrderBelow()")
				return elevio.MD_Down, types.Moving
			}

		}
		log.Println("chooseDir, MD_Up, return MD_STOP")
		return elevio.MD_Stop, types.Idle

	case elevio.MD_Down:
		log.Println("Stop condition: ", e.StopCond)
		if orderBelow(e) {
			if cabOrdersAbove(e) && e.StopCond == types.UpOrder {
				return elevio.MD_Up, types.Moving
			} else {
				log.Println("chooseDir, MD_Down, OrderBelow()")
				return elevio.MD_Down, types.Moving
			}
		}
		if orderAbove(e) {
			if cabOrdersBelow(e) && e.StopCond == types.DownOrder {
				return elevio.MD_Down, types.Moving
			} else {
				log.Println("chooseDir, MD_Down, OrderAbove()")
				return elevio.MD_Up, types.Moving
			}
		}
		log.Println("chooseDir, MD_Down, return MD_STOP")
		return elevio.MD_Stop, types.Idle

	case elevio.MD_Stop:
		log.Println("Stop condition: ", e.StopCond)
		if orderAbove(e) {
			log.Println("chooseDir, MD_Stop, OrderAbove()")
			return elevio.MD_Up, types.Moving
		}
		if orderBelow(e) {
			log.Println("chooseDir, MD_Stop, OrderBelow()")
			return elevio.MD_Down, types.Moving
		}
		log.Println("chooseDir, MD_Stop, return MD_STOP")
		return elevio.MD_Stop, types.Idle

	default:
		//
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

func CheckStopCondition(e *types.Elevator) {
	if hallOrderUpHere(e) {
		e.StopCond = types.UpOrder
	} else if hallOrderDownHere(e) {
		e.StopCond = types.DownOrder
	} else {
		e.StopCond = types.CabOrder
	}
}

func shouldClearAtFloorImmediately(e *types.Elevator, btnFloor int, btnType elevio.ButtonType) bool {
	return e.CurrentFloor == btnFloor &&
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

	log.Println("order direction:", e.OrderDir)
	switch e.OrderDir {
	case elevio.MD_Up:
		if e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallUp] {
			e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallUp] = false
			if IsSingleElevatorMode(e) {
				TurnOffHallLight(elevio.BT_HallUp, e.CurrentFloor)
			}
			log.Println("var en ordre ", elevio.MD_Up, " på floor :", e.CurrentFloor)
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
			log.Println("var en ordre ", elevio.MD_Down, " på floor :", e.CurrentFloor)
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
			log.Println("var en ordre opp i ", elevio.MD_Stop, " på floor :", e.CurrentFloor)
			btnCleared = append(btnCleared, elevio.ButtonEvent{
				Floor:  e.CurrentFloor,
				Button: elevio.BT_HallUp,
			})
		} else if e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallDown] {
			e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallDown] = false
			if IsSingleElevatorMode(e) {
				TurnOffHallLight(elevio.BT_HallDown, e.CurrentFloor)
			}
			log.Println("var en ordre ned i ", elevio.MD_Stop, " på floor :", e.CurrentFloor)
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
	if e.HallOrderMatrix[btnFloor][btnButton] {
		return
	}
	log.Printf("Assigned order -> role:%v floor:%d button:%d\n", ps.Role, btnFloor, btnButton)
	AddOrder(e, btnFloor, btnButton)

	if shouldClearAtFloorImmediately(e, btnFloor, btnButton) {
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
