package elevator

import (
	"Driver-go/elevio"
	Driver "Driver-go/elevio"
	"fmt"
)


func buttonPressedServiceOrder(e *Elevator, btnFloor int, btnType elevio.ButtonType) {
	switch e.state {

	case doorOpen:
		if shouldClearAtFloorImmediately(e, btnFloor, btnType) {
			// må vel ha en timer herr??? dør lys etc ..
		} else {
			addOrder(e, btnFloor, btnType)
		}

	case moving:
		addOrder(e, btnFloor, btnType)

	// hvis vi er i tiolstand idle må vi gjøre noe vi legger til ordre og bruke choose direction til å bestemme hva vi gjør så
	case idle:
		addOrder(e, btnFloor, btnType)

		dir, Nextstate := chooseDirection(e)

		switch Nextstate {
			//hvis chooseDirection sier bli her, stopp og åpne dør
		case doorOpen:
			// motor skal stå stille når døra er åpen
			elevio.SetMotorDirection(elevio.MD_Stop)
			// sett dørlys + start timer 

			clearAtCurrentFloor(e)
			e.state = doorOpen
			e.dir = elevio.MD_Stop
		
			// hvis chooseDirection sier beveg deg sett retnign og tilstand og start motor i riktig retning
		case moving:
			e.state = moving
			e.dir = dir

			elevio.SetMotorDirection(dir)

			// hvis chooseDirection sier ingen ting stopp motor stop heis
		case idle:
			e.state = idle
			e.dir = elevio.MD_Stop
			elevio.SetMotorDirection(elevio.MD_Stop)
		}
	}
}

// hvis vi beveger oss og kommer til en etasje
func serviceOrderAtFloor(e *Elevator, newFloor int) {
	// setter heis etajen til etasjen vi når
	e.currentFloor = newFloor
		
	// denne if setning la jeg til for de at hvis ikke og heisen står stille kan heisen åpne og lukke døra konstant
	if e.state != moving {
        return
    }
	//sjekker om vi skal stopp her
	if shouldStop(e) {
		//hvis ja stopp motor, åpne dør 
		elevio.SetMotorDirection(elevio.MD_Stop)
		e.state = doorOpen
		//sett lys + timer
		clearAtCurrentFloor(e)
			
	
	}
}

