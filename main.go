package main

import (
	"flag"
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/Elevator"
	"heisprosjekt75/MainUtilities"
	"heisprosjekt75/Messages/MessageHandling"
	"heisprosjekt75/Messages/MessageTypes"
	"heisprosjekt75/Network"
	"heisprosjekt75/Network/SendWorldView"
	"heisprosjekt75/Network/bcast"
	"heisprosjekt75/Network/localip"
	"heisprosjekt75/Schedueler"
	"log"
	"time"
)

func main() {
	var elevAddr string
	var stableID string

	const (
		d            = 500 * time.Millisecond
		brodcastPort = 43452
		TCPPort      = "3002"
	)

	flag.StringVar(&stableID, "id", "", "stable elevator id (A/B/C)")
	flag.StringVar(&elevAddr, "elev", "127.0.0.1:15657", "elevator server addr")
	flag.Parse()

	elevio.Init(elevAddr, 4)

	peerID, peerUpdateCh := network.NetworkInit()

	ip, _ := localip.LocalIP()

	e := Elevator.InitNewElevator(peerID, ip)
	e.StableID = stableID

	UDPPrimaryIPIDTx := make(chan sendworldview.PrimaryIPID, 10)
	UDPPrimaryIPIDRx := make(chan sendworldview.PrimaryIPID, 10)
	TCPRx := make(chan messagestypes.Message, 10)
	TCPWorldViewCh := make(chan messagestypes.Message, 10)

	receiveBtnCh := make(chan elevio.ButtonEvent, 10)
	reechFloorCh := make(chan int, 10)
	doorTimeoutCh := make(chan int, 10)
	doorStartTimerCh := make(chan int, 10)
	obstructionBtnCh := make(chan bool)

	go bcast.Transmitter(brodcastPort, UDPPrimaryIPIDTx)
	go bcast.Receiver(brodcastPort, UDPPrimaryIPIDRx)
	go sendworldview.SendPrimaryIpId(UDPPrimaryIPIDTx, d, e)

	go elevio.PollButtons(receiveBtnCh)
	go elevio.PollFloorSensor(reechFloorCh)
	go Elevator.DoorTimeManager(e, doorTimeoutCh, doorStartTimerCh)
	go elevio.PollObstructionSwitch(obstructionBtnCh)
	go Elevator.DoorObstruction(obstructionBtnCh, e, doorStartTimerCh)
	go schedueler.PrimaryMonitorTick(e, doorStartTimerCh, d)

	sendworldview.StartWorldViewSender(&e.Ps, TCPWorldViewCh)

	if elevio.GetFloor() == -1 {
		Elevator.InitBetweenFloor(e)
	}

	for {
		select {
		case btn := <-receiveBtnCh:
			mainutilities.CaseReceiveBtnCh(btn, e, doorStartTimerCh)

		case newFloor := <-reechFloorCh:
			Elevator.FsmServiceOrderAtFloor(e, newFloor, doorStartTimerCh)

		case <-doorTimeoutCh:
			Elevator.DoorTimeout(doorStartTimerCh, e)

		case p := <-peerUpdateCh:
			log.Printf("updating peer ")
			mainutilities.CasePeerUpdateCh(e, doorStartTimerCh, TCPPort, TCPRx, p, TCPWorldViewCh)

		case PrimaryIdIp := <-UDPPrimaryIPIDRx:
			mainutilities.CaseUDPPrimaryIPIDRx(e, TCPPort, TCPRx, PrimaryIdIp)

		case message := <-TCPRx:
			messagelogic.OnMessageReceive(message, e, doorStartTimerCh)
		}
	}
}
