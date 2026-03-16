package main

import (
	"flag"
	"fmt"
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/Elevator"
	"heisprosjekt75/Messages/MessageHandling"
	"heisprosjekt75/Messages/SendMessages"
	"heisprosjekt75/Messages"
	"heisprosjekt75/Network-go/network"
	"heisprosjekt75/Network-go/network/bcast"
	"heisprosjekt75/Network-go/network/localip"
	"heisprosjekt75/Network-go/network/primaryHeartbeat"
	"heisprosjekt75/Network-go/network/tcp"
	"heisprosjekt75/RoleLogic/RoleChanges"
	"heisprosjekt75/RoleLogic/RoleManager"
	"heisprosjekt75/Schedueler"
	"heisprosjekt75/types"
	"log"
	"time"
	"heisprosjekt75/Messages/MessageTypes"
)

func main() {
	var elevAddr string
	var stableID string

	const (
		d            = 500 * time.Millisecond
		brodcastPort = 43452
		TCPPort      = "3001"
	)

	flag.StringVar(&stableID, "id", "", "stable elevator id (A/B/C)")
	flag.StringVar(&elevAddr, "elev", "127.0.0.1:15657", "elevator server addr")
	flag.Parse()

	elevio.Init(elevAddr, 4)

	peerID, peerUpdateCh := network.NetworkInit()

	ip, _ := localip.LocalIP()

	e := Elevator.InitNewElevator(peerID, ip)
	e.StableID = stableID

	UDPHeartbeatTx := make(chan PrimaryHeartbeat.PrimHeartbeat, 10)
	UDPHeartbeatRx := make(chan PrimaryHeartbeat.PrimHeartbeat, 10)
	TCPRx := make(chan messagestypes.Message, 10)
	TCPHeartbeatCh := make(chan messagestypes.Message, 10)

	reaciveBtnCh := make(chan elevio.ButtonEvent, 10)
	reechFloorCh := make(chan int, 10)
	doorTimeoutCh := make(chan int, 10)
	doorStartTimerCh := make(chan int, 10)
	obstructionBtnCh := make(chan bool)

	go bcast.Transmitter(brodcastPort, UDPHeartbeatTx)
	go bcast.Receiver(brodcastPort, UDPHeartbeatRx)
	go PrimaryHeartbeat.SendPrimaryIpId(UDPHeartbeatTx, d, e)

	go elevio.PollButtons(reaciveBtnCh)
	go elevio.PollFloorSensor(reechFloorCh)
	go Elevator.DoorTimeManager(e, doorTimeoutCh, doorStartTimerCh)
	go elevio.PollObstructionSwitch(obstructionBtnCh)
	go Elevator.DoorObstruction(obstructionBtnCh, e, doorStartTimerCh)

	tcp.StartHeartbeatSender(&e.Ps, TCPHeartbeatCh)

	if elevio.GetFloor() == -1 {
		Elevator.InitBetweenFloor(e)
	}

	for {
		select {
		case btn := <-reaciveBtnCh:
			if e.Mode == types.SingleElevator || btn.Button == elevio.BT_Cab {
				Elevator.FsmServiceLocalButton(e, btn.Floor, btn.Button, doorStartTimerCh)
			} else {
				sendmessages.ButtonTransmitLogic(e, btn)
			}

		case newFloor := <-reechFloorCh:
			Elevator.FsmServiceOrderAtFloor(e, newFloor, doorStartTimerCh)

		case <-doorTimeoutCh:
			Elevator.DoorTimeout(doorStartTimerCh, e)

		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %s\n", p.Peers)
			fmt.Printf("  New:      %s\n", p.New)
			fmt.Printf("  Lost:     %s\n", p.Lost)

			RoleManager.RoleElection(p, e, doorStartTimerCh)

			if e.Ps.PrevRole != e.Ps.Role {
				rolechanges.RolesSwitched(TCPPort, TCPRx, e)
				e.Ps.PrevRole = e.Ps.Role

				if e.Ps.Role != types.RolePrimary {
					go tcp.HeartbeatTick(e, 1*time.Second, TCPHeartbeatCh)
					log.Println("entred by stableID:", e.StableID)
				} else {
					go messages.SnapshotTick(e, 500*time.Millisecond)
				}
			}

			rolechanges.HandleLostPeers(p.Lost, e, doorStartTimerCh, p.Peers)

			for stableID, cabs := range types.LostCabOrders {
				log.Println("entred by stableID:", e.StableID)
				log.Printf("lost caborder for stableID: %s, cabs: %v\n", stableID, cabs)
			}

			if len(p.Lost) > 0 && e.Ps.Role == types.RolePrimary && len(p.Peers) > 1 {
				schedueler.PrimarySchedueler(e, doorStartTimerCh)
			}

		case PrimaryIdIp := <-UDPHeartbeatRx:
			oldPrimaryID := e.Ps.PrimaryID

			e.Ps.PrimaryID = PrimaryIdIp.PrimaryID
			e.Ps.PrimaryIP = PrimaryIdIp.PrimaryAddrTCP

			if e.Ps.Role != types.RolePrimary &&
				e.Ps.PrimaryID != "" &&
				e.Ps.PrimaryID != oldPrimaryID {
				log.Println("Primary changed, reconnecting to new primary:", e.Ps.PrimaryID)
				go tcp.ConnectToPrimary(TCPPort, e, TCPRx)
			}

		case message := <-TCPRx:
			messagelogic.OnMessageReceive(message, e, doorStartTimerCh)
		}
	}
}
