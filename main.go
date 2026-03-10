package main

import (
	"flag"
	"fmt"
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/ElevatorP"
	messagelogic "heisprosjekt75/Messages/MessageLogic"
	"heisprosjekt75/Network-go/network"
	"heisprosjekt75/Network-go/network/bcast"
	"heisprosjekt75/Network-go/network/localip"
	PrimaryHeartbeat "heisprosjekt75/Network-go/network/primaryHeartbeat"
	"heisprosjekt75/Network-go/network/tcp"
	rolechanges "heisprosjekt75/RoleLogic/RoleChanges"
	"heisprosjekt75/RoleLogic/RoleManager"

	"heisprosjekt75/types"
	"time"
)

// NB!! Hvis en funksjon ikke funker i main betyr det mest sannsynlig at den er privat, for å dele funkjsoner mellom pakker må forbokstaven være stor
func main() {
	//Initialisering av heiser
	var elevAddr string
	var id string
	const (
		d            = 500 * time.Millisecond
		brodcastPort = 43452
		TCPPort      = "3000"
	)

	//flagene var for å kunne kalle ulike elvator servers i terminalen
	flag.StringVar(&id, "id", "", "node id (A/B/C)")
	flag.StringVar(&elevAddr, "elev", "127.0.0.1:15657", "elevator server addr")
	flag.Parse()

	elevio.Init(elevAddr, 4)
	ps := &types.PeerState{}
	//---------Initialiser nettverk----------------------------------------------------------------------------------------

	// We make channels for sending and receiving our custom data types

	// nettwork init finner noden sin egen id brodacaser herr her jeg og leser om det er andre folk på nettet ved bruk av reive og trancive
	id, peerUpdateCh := network.NetworkInit()
	fmt.Println("min id", id)

	//---------------------------------------------------------------------------------------------------------------------
	ip, _ := localip.LocalIP()

	e := ElevatorP.NewElevator(id, ip)

	UDPHeartbeatTx := make(chan PrimaryHeartbeat.PrimHeartbeat, 10)
	UDPHeartbeatRx := make(chan PrimaryHeartbeat.PrimHeartbeat, 10)
	TCPRx := make(chan tcp.Message, 10)
	TCPHeartbeatCh := make(chan tcp.Message, 10)

	reaciveBtnCh := make(chan elevio.ButtonEvent, 10)
	reechFloorCh := make(chan int, 10)
	doorTimeoutCh := make(chan int, 10)
	doorStartTimerCh := make(chan int, 10)
	obstructionBtnCh := make(chan bool)

	go bcast.Transmitter(brodcastPort, UDPHeartbeatTx)
	go bcast.Receiver(brodcastPort, UDPHeartbeatRx)
	go PrimaryHeartbeat.SendPrimaryIpId(UDPHeartbeatTx, d, ps, e)

	go elevio.PollButtons(reaciveBtnCh)
	go elevio.PollFloorSensor(reechFloorCh)
	go ElevatorP.DoorTimeManager(e, doorTimeoutCh, doorStartTimerCh)
	go elevio.PollObstructionSwitch(obstructionBtnCh)
	go ElevatorP.OnObstruction(obstructionBtnCh, e, doorStartTimerCh)

	tcp.StartHeartbeatSender(ps, TCPHeartbeatCh)

	if elevio.GetFloor() == -1 {
		ElevatorP.OnInitBetweenFloor(e)
	}

	for {
		select {
		case btn := <-reaciveBtnCh:
			if e.Mode == types.SingleElevator || btn.Button == elevio.BT_Cab {
				ElevatorP.ButtonPressedServiceOrder(e, btn.Floor, btn.Button, doorStartTimerCh, ps)
			} else {
				messagelogic.ButtonTransmitLogic(ps, e, btn, doorStartTimerCh)
			}
		case newFloor := <-reechFloorCh:
			ElevatorP.ServiceOrderAtFloor(e, newFloor, doorStartTimerCh, ps)
		case <-doorTimeoutCh:
			ElevatorP.OnDoortimeout(doorStartTimerCh, e)
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
			RoleManager.RoleElection(p, e, ps)
			if ps.PrevRole != ps.Role {
				rolechanges.RolesSwitched(ps, TCPPort, TCPRx, e)
				ps.PrevRole = ps.Role

				if ps.Role == types.RoleBackup {
					go tcp.HeartbeatTick(e, ps, 5*time.Second, TCPHeartbeatCh)
				}
			}

			//if len(p.Lost) > 0 && ps.Role == types.RolePrimary {
			//	schedueler.MasterSchedueler(e, ps, doorStartTimerCh)
			//	fmt.Printf("Primary lost, redeligating orders \n")
			//}

			// TEST: hvis jeg er primary og har backupID

		case PrimaryIdIp := <-UDPHeartbeatRx:
			//fmt.Printf("  PrimaryID:    %q\n", PrimaryIdIp.PrimaryID)
			//fmt.Printf("  PrimaryIP:    %q\n", PrimaryIdIp.PrimaryAddrTCP)

			ps.PrimaryID = PrimaryIdIp.PrimaryID
			ps.PrimaryIP = PrimaryIdIp.PrimaryAddrTCP
			//btn := elevio.ButtonEvent{Floor: 2, Button: elevio.BT_HallDown}
			//schedueler.DelegateOrders(ps.BackupID, ps, e, btn)

		case message := <-TCPRx:

			messagelogic.OnMessageReceive(message, ps, e, doorStartTimerCh)

		}
		// til senere.....

		//trenger vel melding inn melding ut kanal?
		// og en assigned order, stateheartbeat kanal?
		// go routine for recive and send

	}
}
