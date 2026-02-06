package ElevatorP

import (
	"Driver-go/elevio"
	"fmt"
	
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

// sjekker er det en caborder
func cabOrdersHere(e *Elevator) bool {
	return e.CabOrderMatrix[e.CurrentFloor][0]
}
// sjekker er det en hallOrder Opp
func hallOrderUpHere(e *Elevator) bool {
	return e.HallorderMatrix[e.CurrentFloor][elevio.BT_HallUp]
}
// sjekker er det en hallOrder ned
func hallOrderDownHere(e *Elevator) bool {
	return e.HallorderMatrix[e.CurrentFloor][elevio.BT_HallDown]
}

// bare lagde en ekstra funskjonm når vi skulle sjekke eksisterere det en ordre her
func ordersHere(e *Elevator) bool {
	return cabOrdersHere(e) || hallOrderUpHere(e) || hallOrderDownHere(e)
}

// sjekker om det er ordre
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

// sjekke rom det er orde under etasjen vi er i
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


//logikk for å velge retning
func chooseDirection(e *Elevator) (elevio.MotorDirection, elevatorState) {
	fmt.Print("We are now in choosDirection func, Moving to switch case \n")
	switch e.Dir{

	case elevio.MD_Up:
		fmt.Print("Elevator is Moving up \n")
// hvis retning er opp, sjekke om det finnes noen ordre her da vil jeg stoppe og tilstand blir dør åpen
		if ordersHere(e) {
			return elevio.MD_Stop, DoorOpen
		}
// hvis retning er opp, sjekke om det finnes noen ordre over meg så vil jeg fortsette oppover
		if orderAbove(e) {
			return elevio.MD_Up, Moving
		}

// hvis retning er opp, sjekke om det finnes noen ordre under meg da vil jeg snu retning
		if orderBelow(e) {
			return elevio.MD_Down, Moving
		}
		return elevio.MD_Stop, Idle
// samme logikk for de andre retningene
	case elevio.MD_Down:
		fmt.Print("Elevator Moving down \n")
		if ordersHere(e) {
			return elevio.MD_Stop, DoorOpen
		}
		if orderBelow(e) {
			return elevio.MD_Down, Moving
		}
		if orderAbove(e) {
			return elevio.MD_Up, Moving
		}
		return elevio.MD_Stop, Idle

	case elevio.MD_Stop:
		fmt.Print("Elevator not Moving\n")
		if ordersHere(e) {
			return elevio.MD_Stop, DoorOpen
		}
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

// logikk for når heisen skal stoppe
func shouldStop(e *Elevator) bool {
	switch e.Dir{
// hvis retning er opp sjekk finnes det en cab order eller halorderopp stopp, eller hvis det ikke finnes noen ordre over meg så stopp
	case elevio.MD_Up:
		return cabOrdersHere(e) || hallOrderUpHere(e) || !orderAbove(e)
// hvis retning er ned sjekk finnes det en cab order eller halorder ned stopp, eller hvis det ikke finnes noen ordre under meg meg så stopp
	case elevio.MD_Down:
		return cabOrdersHere(e) || hallOrderDownHere(e) || !orderBelow(e)

	case elevio.MD_Stop:
		return true

	default:
		return true
	}
}

// når en skal ta en ordre bare med en gang, hvis etasjen vi er i  er likk knappe trykket og retningen er enten opp og knappen var opp eller retning ned kanpp var ned eller retning er stopp eller knappen var en cab
func shouldClearAtFloorImmediately(e *Elevator, btnFloor int, btnType elevio.ButtonType) bool {
	return e.CurrentFloor == btnFloor &&
		((e.Dir == elevio.MD_Up && btnType == elevio.BT_HallUp) ||
			(e.Dir == elevio.MD_Down && btnType == elevio.BT_HallDown) ||
			(e.Dir == elevio.MD_Stop) ||
			(btnType == elevio.BT_Cab))
}

// logikk for å fjerne ordre
func clearAtCurrentFloor(e *Elevator, prevDir elevio.MotorDirection) {
	
	//vi fjerner alltid hvis det var en cab order
	e.CabOrderMatrix[e.CurrentFloor][0] = false
	TurnOffCabLight(e.CurrentFloor)

	switch prevDir{

	case elevio.MD_Up:
		//hvis retnigne er opp fjerner en hallup orderen
		e.HallorderMatrix[e.CurrentFloor][elevio.BT_HallUp] = false
		TurnOffHallLight(elevio.BT_HallUp, e.CurrentFloor)
		//og hvis ingen ordre over meg og og ingen hallorder i etasjen jeg er i kan en ekspeder ned ordre
		if !orderAbove(e) && !hallOrderUpHere(e) {
			e.HallorderMatrix[e.CurrentFloor][elevio.BT_HallDown] = false
			TurnOffHallLight(elevio.BT_HallDown, e.CurrentFloor)
		}
// samme logikk hvis retning er er ned
	case elevio.MD_Down:
		
		e.HallorderMatrix[e.CurrentFloor][elevio.BT_HallDown] = false
		TurnOffHallLight(elevio.BT_HallDown, e.CurrentFloor)
		
		if !orderBelow(e) && !hallOrderDownHere(e) {
			e.HallorderMatrix[e.CurrentFloor][elevio.BT_HallUp] = false
			TurnOffHallLight(elevio.BT_HallUp, e.CurrentFloor)
		}

	case elevio.MD_Stop:
		// logikk hvor å vite hvilken ordre heisen skal ta her
	}
}
