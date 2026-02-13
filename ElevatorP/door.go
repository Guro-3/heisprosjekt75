package ElevatorP

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

func onDoorOpen(doorStartTimerCh chan int, e *Elevator) {
	prevDir := e.Dir
	e.State = DoorOpen

	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetDoorOpenLamp(true)

	fmt.Print("Door open \n")
	clearAtCurrentFloor(e, prevDir)
	doorStartTimerCh <- timeDoorOpenDuration
	// trenger nok et bedre funksjons navn, men skal sette på lys og dør timer, og choose nex direction
	// sender til kanal start timer
}

func OnDoortimeout(doorStartTimerCh chan int, e *Elevator) {
	fmt.Print("Doors closing \n")
	if e.obstructed {
		fmt.Print("timout ignored\n")
		return
	}

	fmt.Print("Door closing")
	elevio.SetDoorOpenLamp(false)
	StartAction(e)
	// hva gjør vi når tiden er ute
}

func OnObstruction(obstructionBtnCh chan bool, e *Elevator, doorStartTimerCh chan int) {
	for {

		obstruction := <-obstructionBtnCh

		if obstruction {
			fmt.Printf("elevator state:%d\n", e.State)
			e.obstructed = true

			if e.State == DoorOpen || e.State == Idle {
				fmt.Print(" door obstruction \n")

				elevio.SetMotorDirection(elevio.MD_Stop)
				elevio.SetDoorOpenLamp(true)
			}

		} else {
			fmt.Print("obstruction cleared\n")
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
			fmt.Print("Timer has started\n")
			select {
			case <-timer.C:
				fmt.Print("time out\n")
				doorTimeoutCh <- timeDuration
			}
		}
	}

}
