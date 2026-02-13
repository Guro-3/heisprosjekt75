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

	switch e.State {

	case DoorOpen:
		if shouldClearAtFloorImmediately(e, btnFloor, btnType) {
			onDoorOpen(doorStartTimerCh, e)
			TurnOffHallLight(btnType, btnFloor)

		} else {
			addOrder(e, btnFloor, btnType)
		}

	case Moving:

		addOrder(e, btnFloor, btnType)

	case Idle:

		if shouldClearAtFloorImmediately(e, btnFloor, btnType) {
			onDoorOpen(doorStartTimerCh, e)
			TurnOffHallLight(btnType, btnFloor)
			return
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

	switch Nextstate {
	case Moving:
		e.State = Moving
		e.Dir = Dir
		elevio.SetMotorDirection(Dir)

	case Idle:
		e.State = Idle
		e.Dir = elevio.MD_Stop
		elevio.SetMotorDirection(elevio.MD_Stop)
	}
}


func ServiceOrderAtFloor(e *Elevator, newFloor int, doorStartTimerCh chan int) {
	e.CurrentFloor = newFloor
	FloorLight(e)

	if e.State != Moving {
		return
	}
	if shouldStop(e) {
		onDoorOpen(doorStartTimerCh, e)
	}
}

func OnInitBetweenFloor(e *Elevator) {
	elevio.SetMotorDirection(elevio.MD_Down)
	e.Dir = elevio.MD_Down
	e.State = Moving
}
