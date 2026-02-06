package ElevatorP

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)


func onDoorOpen(doorStartTimerCh chan int, e *Elevator){
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

func OnDoortimeout(doorStartTimerCh chan int, e *Elevator){
	fmt.Print("Doors closing \n")

	elevio.SetDoorOpenLamp(false)
	StartAction(doorStartTimerCh, e)
	// hva gjør vi når tiden er ute
}

func DoorTimeManager(doorTimeoutCh chan int, doorStartTimerCh chan int){
	
	for{
		select{
		case timeDuration := <- doorStartTimerCh :
			timer := time.NewTimer(time.Duration(timeDuration)*time.Second)
			fmt.Print("Timer has started\n")
			select{
				//legg in obstruction
			case <- timer.C:
				fmt.Print("time out\n")
				doorTimeoutCh <- timeDuration 
			}
		}
	}

}

			
	