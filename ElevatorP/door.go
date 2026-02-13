package ElevatorP

import (
	"Driver-go/elevio"
	"time"
)

func onDoorOpen(doorStartTimerCh chan int, e *Elevator) {
	prevDir := e.Dir
	e.State = DoorOpen

	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetDoorOpenLamp(true)

	clearAtCurrentFloor(e, prevDir)
	doorStartTimerCh <- timeDoorOpenDuration
}

func OnDoortimeout(doorStartTimerCh chan int, e *Elevator) {
	if e.obstructed {
		return
	}
	elevio.SetDoorOpenLamp(false)
	StartAction(e)
}

func OnObstruction(obstructionBtnCh chan bool, e *Elevator, doorStartTimerCh chan int) {
	for {
		obstruction := <-obstructionBtnCh

		if obstruction {
			e.obstructed = true
			if e.State == DoorOpen || e.State == Idle {
				elevio.SetMotorDirection(elevio.MD_Stop)
				elevio.SetDoorOpenLamp(true)
			}
		} else {			
			e.obstructed = false
			doorStartTimerCh <- timeDoorOpenDuration
		}
	}
}

func DoorTimeManager(e *Elevator, doorTimeoutCh chan int, doorStartTimerCh chan int) {
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
