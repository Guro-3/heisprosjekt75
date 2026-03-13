package ElevatorP

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/types"
	"time"
)

func onDoorOpen(doorStartTimerCh chan int, e *types.Elevator, ps *types.PeerState) {
	prevDir := e.Dir
	e.State = types.DoorOpen

	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetDoorOpenLamp(true)

	clearAtCurrentFloor(e, prevDir,ps)
	doorStartTimerCh <- types.TimeDoorOpenDuration
}

func OnDoortimeout(doorStartTimerCh chan int, e *types.Elevator) {
	if e.Obstructed {
		return
	}
	elevio.SetDoorOpenLamp(false)
	StartAction(e)
}

func OnObstruction(obstructionBtnCh chan bool, e *types.Elevator, doorStartTimerCh chan int) {
	for {
		obstruction := <-obstructionBtnCh

		if obstruction {
			e.Obstructed = true
			if e.State == types.DoorOpen || e.State == types.Idle {
				elevio.SetMotorDirection(elevio.MD_Stop)
				elevio.SetDoorOpenLamp(true)
			}
		} else {
			e.Obstructed = false
			doorStartTimerCh <- types.TimeDoorOpenDuration
		}
	}
}

func DoorTimeManager(e *types.Elevator, doorTimeoutCh chan int, doorStartTimerCh chan int) {
	for {
		select {
		case timeDuration := <-doorStartTimerCh:
			timer := time.NewTimer(time.Duration(timeDuration) * time.Second)
			select {
			case <-timer.C:
				doorTimeoutCh <- timeDuration
			}
		}
	}
}
