package main

import (
	"Driver-go/elevio"
	"github.com/Guro-3/heisprosjekt75.git/ElevatorP"
	"Network-go/network"
	"Network-go/network/bcast"
	"Network-go/network/localip"
	"Network-go/network/peers"
	"fmt"
)


// NB!! Hvis en funksjon ikke funker i main betyr det mest sannsynlig at den er privat, for å dele funkjsoner mellom pakker må forbokstaven være stor
func main() {
	//Initialisering av heiser
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

//---------Initialiser nettverk----------------------------------------------------------------------------------------
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We make channels for sending and receiving our custom data types
	UDPHeartbeatTx := make(chan ElevatorP.Heartbeat)
	UDPHeartbeatRx := make(chan ElevatorP.Heartbeat)
	
	network.NetworkInit()
//---------------------------------------------------------------------------------------------------------------------
	for {
		select {
		case btn := <-reaciveBtnCh:
			ElevatorP.ButtonPressedServiceOrder(e, btn.Floor, btn.Button, doorStartTimerCh)
		case newFloor := <-reechFloorCh:
			ElevatorP.ServiceOrderAtFloor(e, newFloor, doorStartTimerCh)
		case <-doorTimeoutCh:
			ElevatorP.OnDoortimeout(doorStartTimerCh, e)
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-UDPHeartbeatRx:
			fmt.Printf("Received: %#v\n", a)
		}
	}
	// til senere.....

	//trenger vel melding inn melding ut kanal?
	// og en assigned order, stateheartbeat kanal?
	// go routine for recive and send

}
