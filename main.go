package main

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/ElevatorP"
	"heisprosjekt75/Network-go/network"

	//"heisprosjekt75/RoleManager"

	//"Network-go/network/localip"
	// "Network-go/network/peers"
	"flag"
	"fmt"
)

// NB!! Hvis en funksjon ikke funker i main betyr det mest sannsynlig at den er privat, for å dele funkjsoner mellom pakker må forbokstaven være stor
func main() {
	//Initialisering av heiser
	var elevAddr string
	var id string

	//flagene var for å kunne kalle ulike elvator servers i terminalen
	flag.StringVar(&id, "id", "", "node id (A/B/C)")
	flag.StringVar(&elevAddr, "elev", "127.0.0.1:15657", "elevator server addr")
	flag.Parse()

	elevio.Init(elevAddr, 4)
	//ps := &RoleManager.PeerState{}
	//---------Initialiser nettverk----------------------------------------------------------------------------------------

	// We make channels for sending and receiving our custom data types
	// UDPHeartbeatTx := make(chan ElevatorP.Heartbeat)
	// UDPHeartbeatRx := make(chan ElevatorP.Heartbeat)
	// nettwork init finner noden sin egen id brodacaser herr her jeg og leser om det er andre folk på nettet ved bruk av reive og trancive
	id, peerUpdateCh := network.NetworkInit()
	fmt.Printf("min id%d\n", id)

	//---------------------------------------------------------------------------------------------------------------------

	e := ElevatorP.NewElevator(id)
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
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
			//RoleManager.RoleElection(p, e.MyID, ps)

			// case a := <-UDPHeartbeatRx:
			// 	fmt.Printf("Received: %#v\n", a)
			// }
		}
		// til senere.....

		//trenger vel melding inn melding ut kanal?
		// og en assigned order, stateheartbeat kanal?
		// go routine for recive and send

	}
}
