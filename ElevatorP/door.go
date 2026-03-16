package ElevatorP

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/types"
	"time"
)

func onDoorOpen(doorStartTimerCh chan int, e *types.Elevator, ps *types.PeerState) {
	if elevio.GetFloor() == -1 {
		return
	}
	e.State = types.DoorOpen

	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetDoorOpenLamp(true)

	clearAtCurrentFloor(e, e.OrderDir, ps)

	e.ClearedRevDir = shouldClearOppositeAtCurrentFloor(e)

	doorStartTimerCh <- types.TimeDoorOpenDuration
}

func OnDoortimeout(doorStartTimerCh chan int, e *types.Elevator, ps *types.PeerState) {
	if e.Obstructed {
		return
	}
	if e.ClearedRevDir {
		clearOppositeAtAtCurrentFloor(e, ps)

		switch e.OrderDir {
		case elevio.MD_Up:
			e.OrderDir = elevio.MD_Down
		case elevio.MD_Down:
			e.OrderDir = elevio.MD_Up
		}

		e.ClearedRevDir = false
		doorStartTimerCh <- types.TimeDoorOpenDuration
		return
	}

	elevio.SetDoorOpenLamp(false)
	StartAction(e, doorStartTimerCh, ps)
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
			if e.State == types.DoorOpen {
				doorStartTimerCh <- types.TimeDoorOpenDuration
			}
		}
	}
}

func DoorTimeManager(e *types.Elevator, doorTimeoutCh chan int, doorStartTimerCh chan int) {
	var timer *time.Timer
	var timerChan <-chan time.Time
	for {
		select {
		case timeDuration := <-doorStartTimerCh:
			if timer != nil {
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
			}
			timer = time.NewTimer(time.Duration(timeDuration) * time.Second)
			timerChan = timer.C
		case <-timerChan:
			timerChan = nil
			if !e.Obstructed {
				doorTimeoutCh <- types.TimeDoorOpenDuration
			}

		}
	}
}
