package ElevatorP

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/types"
)



func ButtonPressedServiceOrder(e *types.Elevator, btnFloor int, btnType elevio.ButtonType, doorStartTimerCh chan int, ps *types.PeerState) {
	
	switch e.State {

	case types.DoorOpen:
		if shouldClearAtFloorImmediately(e, btnFloor, btnType) {
			TurnOffHallLight(btnType, btnFloor)
			onDoorOpen(doorStartTimerCh, e, ps)

		} else {
			AddOrder(e, btnFloor, btnType)
		}

	case types.Moving:

		AddOrder(e, btnFloor, btnType)

	case types.Idle:

		if shouldClearAtFloorImmediately(e, btnFloor, btnType) {
			TurnOffHallLight(btnType, btnFloor)
			onDoorOpen(doorStartTimerCh, e, ps)
			return
		}

		AddOrder(e, btnFloor, btnType)
		StartAction(e)
	}
}

func StartAction(e *types.Elevator) {
	if e.Obstructed {
		return
	}
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


func ServiceOrderAtFloor(e *types.Elevator, newFloor int, doorStartTimerCh chan int, ps *types.PeerState) {
	e.CurrentFloor = newFloor
	FloorLight(e)

	if e.State != types.Moving {
		return
	}
	if shouldStop(e) {
		onDoorOpen(doorStartTimerCh, e, ps)
	}
}

func OnInitBetweenFloor(e *types.Elevator) {
	elevio.SetMotorDirection(elevio.MD_Down)
	e.Dir = elevio.MD_Down
	e.State = types.Moving
}
