package Elevator

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/types"
)

func FsmServiceLocalButton(e *types.Elevator, btnFloor int, btnType elevio.ButtonType, doorStartTimerCh chan int) {

	switch e.State {

	case types.DoorOpen:
		if shouldClearOrderAtFloorImmediately(e, btnFloor, btnType) {
			DoorOpen(doorStartTimerCh, e)

		} else {
			AddOrder(e, btnFloor, btnType)
		}

	case types.Moving:

		AddOrder(e, btnFloor, btnType)

	case types.Idle:

		if shouldClearOrderAtFloorImmediately(e, btnFloor, btnType) {
			DoorOpen(doorStartTimerCh, e)
			return
		}

		AddOrder(e, btnFloor, btnType)
		FsmStartAction(e, doorStartTimerCh)
	}
}

func FsmStartAction(e *types.Elevator, doorStartTimerCh chan int) {
	if e.Obstructed {
		return
	}
	if e.State == types.Moving {
		return
	}
	/*for {
		if e.State == types.DoorOpen {
			continue
		} else {
			break
		}
	}*/

	Dir, Nextstate := chooseDirection(e)

	switch Nextstate {
	case types.Moving:
		e.State = types.Moving
		e.Dir = Dir
		elevio.SetMotorDirection(Dir)

	case types.Idle:
		e.State = types.Idle
		e.Dir = elevio.MD_Stop
		elevio.SetMotorDirection(elevio.MD_Stop)
	}
}

func FsmServiceOrderAtFloor(e *types.Elevator, newFloor int, doorStartTimerCh chan int) {
	e.CurrentFloor = newFloor
	FloorLight(e)

	if e.Initializing {
		elevio.SetMotorDirection(elevio.MD_Stop)
		e.State = types.Idle
		e.Dir = elevio.MD_Stop
		e.Initializing = false
		return
	}
	if e.State != types.Moving {
		return
	}
	if shouldStop(e) {
		setOrderDirAtStop(e)
		DoorOpen(doorStartTimerCh, e)
	}
}
