package main

import (
	"Driver-go/elevio"
	"fmt"

	"github.com/Guro-3/heisprosjekt75.git/ElevatorP"
)

// NB!! Hvis en funksjon ikke funker i main betyr det mest sannsynlig at den er privat, for å dele funkjsoner mellom pakker må forbokstaven være stor
func main() {
	// 1. initialisere server
	elevio.Init("127.0.0.1", 4)
	//2. lage kanaler som go rutinene kan bruke
	e := ElevatorP.NewElevator()
	reaciveBtnCh := make(chan elevio.ButtonEvent, 10)
	reechFloorCh := make(chan int, 10)
	doorTimeoutCh := make(chan int)
	doorStartTimerCh := make(chan int)

	go elevio.PollButtons(reaciveBtnCh)
	go elevio.PollFloorSensor(reechFloorCh)
	go ElevatorP.DoorTimeManager(doorTimeoutCh, doorStartTimerCh)

	//en dør go funksjon, som starter timer sedenr door timou tilbake
	// den må få et start dør event, starte ny timer med ønsket dyration
	// tid ute sender ut på kanal timout
	if elevio.GetFloor() == -1 {
		ElevatorP.OnInitBetweenFloor(e)
	}

	fmt.Print("Starter fsm for loop")
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
	// 4. kjører GoToNearestFloor elller noe i den duren
	// 5. starte er for select loop:
	// hvis det kommer noe på kanal button: kall ButtonPressedServiceOrder
	//hvis det kommer noe på kanal floor reached: kall serviceOrderAtFloor
	// hvis det kommer noe på kanal start timer blir start dør timer kalt, og det kommer noe på stopptimer blir doortimout kaldt

	// til senere.....
	//trenger vel melding inn melding ut kanal?
	// og en assigned order, stateheartbeat kanal?

	// go routine for recive and send

}
