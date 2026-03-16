package Elevator

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/types"
	"time"
)

func DoorOpen(doorStartTimerCh chan int, e *types.Elevator) {
	if elevio.GetFloor() == -1 {
		return
	}
	e.State = types.DoorOpen

	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetDoorOpenLamp(true)

	clearOrderAtCurrentFloor(e, e.OrderDir)

	e.ClearedRevDir = shouldClearOppositeOrderAtCurrentFloor(e)

	doorStartTimerCh <- types.TimeDoorOpenDuration
}

func DoorTimeout(doorStartTimerCh chan int, e *types.Elevator) {
	if e.Obstructed {
		return
	}
	if e.ClearedRevDir {
		clearOppositeOrderAtCurrentFloor(e)

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
	FsmStartAction(e, doorStartTimerCh)
}

func DoorObstruction(obstructionBtnCh chan bool, e *types.Elevator, doorStartTimerCh chan int) {
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
