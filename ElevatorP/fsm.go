package ElevatorP

import (
	"Driver-go/elevio"
	"fmt"
)

func addOrderLocal(e *Elevator) bool {
	return e.Mode == SingleElavator
}

func ButtonPressedServiceOrder(e *Elevator, btnFloor int, btnType elevio.ButtonType, doorStartTimerCh chan int) {
	fmt.Print("In func buttonPressedServiceOrder: \n")
	if !addOrderLocal(e) {
		fmt.Print("Multiple elevators online \n")
		/// gjør noe
		return
	}
	fmt.Print("Elevator is in single Mode \n")
	fmt.Print("we are now Moving on to switch case for singles:) \n")

	switch e.State {

	case DoorOpen:
		fmt.Print("Door open \n")
		if shouldClearAtFloorImmediately(e, btnFloor, btnType) {
			onDoorOpen(doorStartTimerCh, e)
			TurnOffHallLight(btnType, btnFloor) // kan ikke lyse være cab her?
			// må ha dør lys etc ..
		} else {
			addOrder(e, btnFloor, btnType)
		}

	case Moving:
		fmt.Print("Moving and adding order\n")
		addOrder(e, btnFloor, btnType)

	// hvis vi er i tilstand Idle må vi gjøre noe vi legger til ordre og bruke choose direction til å bestemme hva vi gjør så
	case Idle:
		fmt.Print("in Idle\n")
		if shouldClearAtFloorImmediately(e, btnFloor, btnType) {
			onDoorOpen(doorStartTimerCh, e)
			TurnOffHallLight(btnType, btnFloor)
			return // burde vi her også sjekke shouldClearImedetly?
		}

		addOrder(e, btnFloor, btnType)
		StartAction(e)
	}
}

func StartAction(e *Elevator) {
	if e.obstructed {
		return
	}
	Dir, Nextstate := chooseDirection(e)
	fmt.Print("Switch case in Idle\n")

	switch Nextstate {
	// hvis chooseDirection sier bli her, stopp og åpne dør
	// case DoorOpen:
	// 	onDoorOpen(doorStartTimerCh, e)
	// 	// motor skal stå stille når døra er åpen
	// 	// sett dørlys + start timer
	// 	e.State = DoorOpen
	// 	e.Dir = elevio.MD_Stop
	//
	// markerte ut denne logikken fordi ettet tenking burde ikke open dår være noe nexstate gir?

	// hvis chooseDirection sier beveg deg sett retnign og tilstand og start motor i riktig retning
	case Moving:
		e.State = Moving
		e.Dir = Dir
		elevio.SetMotorDirection(Dir)

	// hvis chooseDirection sier ingen ting stopp motor stop heis
	case Idle:
		e.State = Idle
		e.Dir = elevio.MD_Stop
		elevio.SetMotorDirection(elevio.MD_Stop)
	}
}

// hvis vi beveger oss og kommer til en etasje
func ServiceOrderAtFloor(e *Elevator, newFloor int, doorStartTimerCh chan int) {
	// setter heis etajen til etasjen vi når
	e.CurrentFloor = newFloor
	FloorLight(e)

	// denne if setning la jeg til for de at hvis ikke og heisen står stille kan heisen åpne og lukke døra konstant
	if e.State != Moving {
		return
	}

	//sjekker om vi skal stopp her
	if shouldStop(e) {
		//hvis ja stopp motor, åpne dør
		onDoorOpen(doorStartTimerCh, e)
		//sett lys + timer
	}
}

func OnInitBetweenFloor(e *Elevator) {
	fmt.Print("Init for elevator between floors\n")
	elevio.SetMotorDirection(elevio.MD_Down)
	e.Dir = elevio.MD_Down
	e.State = Moving
}
