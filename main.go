package main

import (
	"Driver-go/elevio"
	"github.com/Guro-3/heisprosjekt75.git/ElevatorP"
)

// NB!! Hvis en funksjon ikke funker i main betyr det mest sannsynlig at den er privat, for å dele funkjsoner mellom pakker må forbokstaven være stor
func main() {
	elevio.Init("127.0.0.1:15657", 4)

	e := ElevatorP.NewElevator()
	reaciveBtnCh := make(chan elevio.ButtonEvent, 10)
	reechFloorCh := make(chan int, 10)
	doorTimeoutCh := make(chan int, 10)
	doorStartTimerCh := make(chan int, 10)
	obstructionBtnCh := make(chan bool)

	go elevio.PollButtons(reaciveBtnCh)
	go elevio.PollFloorSensor(reechFloorCh)
	go ElevatorP.DoorTimeManager(e, doorTimeoutCh, doorStartTimerCh)
	go elevio.PollObstructionSwitch(obstructionBtnCh)
	go ElevatorP.OnObstruction(obstructionBtnCh, e, doorStartTimerCh)

	if elevio.GetFloor() == -1 {
		ElevatorP.OnInitBetweenFloor(e)
	}

	
	for {
		select {
		case btn := <-reaciveBtnCh:
			ElevatorP.ButtonPressedServiceOrder(e, btn.Floor, btn.Button, doorStartTimerCh)
		case newFloor := <-reechFloorCh:
			ElevatorP.ServiceOrderAtFloor(e, newFloor, doorStartTimerCh)
		case <-doorTimeoutCh:
			ElevatorP.OnDoortimeout(doorStartTimerCh, e)

		}
	}
	// til senere.....

	//trenger vel melding inn melding ut kanal?
	// og en assigned order, stateheartbeat kanal?
	// go routine for recive and send

}
