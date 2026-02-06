package elevator

import (
	"Driver-go/elevio"
	"fmt"
)

func initreturntofloor()

func addOrderLocal(e *Elevator) bool{
	return e.mode == singleElavator
}



func buttonPressedServiceOrder(e *Elevator, btnFloor int, btnType elevio.ButtonType) {
	fmt.Print("In func buttonPressedServiceOrder: \n")
	if !addOrderLocal(e){
		fmt.Print("Multiple elevators online \n")
		/// gjør noe
		return
	}
	fmt.Print("Elevator is in single mode \n")
	fmt.Print("we are now moving on to switch case for singles:) \n")
	switch e.state {

	case doorOpen:
		fmt.Print("Door open \n")
		if shouldClearAtFloorImmediately(e, btnFloor, btnType) {
			// må vel ha en timer her??? dør lys etc ..
		} else {
			addOrder(e, btnFloor, btnType)
		}
		
	case moving:
		fmt.Print("Moving and adding order\n")
		addOrder(e, btnFloor, btnType)

	// hvis vi er i tilstand idle må vi gjøre noe vi legger til ordre og bruke choose direction til å bestemme hva vi gjør så
	case idle:
		fmt.Print("in Idle\n")
		addOrder(e, btnFloor, btnType)

		dir, Nextstate := chooseDirection(e)
		fmt.Print("Switch case in Idle\n")
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
	
func onDoorOpen(){
	// trenger nok et bedre funksjons navn, men skal sette på lys og dør timer, og choose nex direction
	// sender til kanal start timer
}
func onDoortimeout(){
	// hva gjør vi når tiden er ute
}

