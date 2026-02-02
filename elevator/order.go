package elevator

import "Driver-go/elevio"


func addOrder()

// sjekker er det en caborder
func cabOrdersHere(e *Elevator) bool {
	return e.cabOrderMatrix[e.currentFloor][0]
}
// sjekker er det en hallOrder Opp
func hallOrderUpHere(e *Elevator) bool {
	return e.HallorderMatrix[e.currentFloor][elevio.BT_HallUp]
}
// sjekker er det en hallOrder ned
func hallOrderDownHere(e *Elevator) bool {
	return e.HallorderMatrix[e.currentFloor][elevio.BT_HallDown]
}

// bare lagde en ekstra funskjonm når vi skulle sjekke eksisterere det en ordre her
func ordersHere(e *Elevator) bool {
	return cabOrdersHere(e) || hallOrderUpHere(e) || hallOrderDownHere(e)
}

// sjekker om det er ordre
func orderBelow(e *Elevator) bool {
	for f := e.currentFloor - 1; f >= 0; f-- {
		for b := 0; b < numHallButtons; b++ {
			if e.HallorderMatrix[f][b] {
				return true
			}
		}
		
		if e.cabOrderMatrix[f][0] {
			return true
		}
	}
	return false
}

// sjekke rom det er orde under etasjen vi er i
func orderAbove(e *Elevator) bool {
	for f := e.currentFloor + 1; f < NUMFloors; f++ {
		for b := 0; b < numHallButtons; b++ {
			if e.HallorderMatrix[f][b] {
				return true
			}
		}
		if e.cabOrderMatrix[f][0] {
			return true
		}
	}
	return false
}


//logikk for å velge retning
func chooseDirection(e *Elevator) (elevio.MotorDirection, elevatorState) {
	switch e.dir {
// hvis retning er opp, sjekke om det finnes noen ordre over meg så vil jeg fortsette oppover
	case elevio.MD_Up:
		if orderAbove(e) {
			return elevio.MD_Up, moving
		}
// hvis retning er opp, sjekke om det finnes noen ordre her da vil jeg stoppe og tilstand blir dør åpen
		if ordersHere(e) {
			return elevio.MD_Stop, doorOpen
		}
// hvis retning er opp, sjekke om det finnes noen ordre under meg da vil jeg snu retning
		if orderBelow(e) {
			return elevio.MD_Down, moving
		}
		return elevio.MD_Stop, idle
// samme logikk for de andre retningene
	case elevio.MD_Down:
		if orderBelow(e) {
			return elevio.MD_Down, moving
		}
		if ordersHere(e) {
			return elevio.MD_Stop, doorOpen
		}
		if orderAbove(e) {
			return elevio.MD_Up, moving
		}
		return elevio.MD_Stop, idle

	case elevio.MD_Stop:
		if ordersHere(e) {
			return elevio.MD_Stop, doorOpen
		}
		if orderAbove(e) {
			return elevio.MD_Up, moving
		}
		if orderBelow(e) {
			return elevio.MD_Down, moving
		}
		return elevio.MD_Stop, idle

	default:
		return elevio.MD_Stop, idle
	}
}

// logikk for når heisen skal stoppe
func shouldStop(e *Elevator) bool {
	switch e.dir {
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
	return e.currentFloor == btnFloor &&
		((e.dir == elevio.MD_Up && btnType == elevio.BT_HallUp) ||
			(e.dir == elevio.MD_Down && btnType == elevio.BT_HallDown) ||
			(e.dir == elevio.MD_Stop) ||
			(btnType == elevio.BT_Cab))
}

// logikk for å fjerne ordre
func clearAtCurrentFloor(e *Elevator) {
	
	//vi fjerner alltid hvis det var en cab order
	e.cabOrderMatrix[e.currentFloor][0] = false

	switch e.dir {

	case elevio.MD_Up:
		//hvis retnigne er opp fjerner en hallup orderen
		e.HallorderMatrix[e.currentFloor][elevio.BT_HallUp] = false

		//og hvis ingen ordre over meg og og ingen hallorder i etasjen jeg er i kan en ekspeder ned ordre
		if !orderAbove(e) && !hallOrderUpHere(e) {
			e.HallorderMatrix[e.currentFloor][elevio.BT_HallDown] = false
		}
// samme logikk hvis retning er er ned
	case elevio.MD_Down:
		
		e.HallorderMatrix[e.currentFloor][elevio.BT_HallDown] = false

		
		if !orderBelow(e) && !hallOrderDownHere(e) {
			e.HallorderMatrix[e.currentFloor][elevio.BT_HallUp] = false
		}

	case elevio.MD_Stop:
		// logikk hvor å vite hvilken ordre heisen skal ta her
	}
}
