package main

import (
	"flag"
	"fmt"
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/ElevatorP"
	messagelogic "heisprosjekt75/Messages/MessageLogic"
	sendmessages "heisprosjekt75/Messages/SendMessages"
	"heisprosjekt75/Network-go/network"
	"heisprosjekt75/Network-go/network/bcast"
	"heisprosjekt75/Network-go/network/localip"
	PrimaryHeartbeat "heisprosjekt75/Network-go/network/primaryHeartbeat"
	"heisprosjekt75/Network-go/network/tcp"
	rolechanges "heisprosjekt75/RoleLogic/RoleChanges"
	"heisprosjekt75/RoleLogic/RoleManager"
	schedueler "heisprosjekt75/Schedueler"
	"heisprosjekt75/types"
	"log"
	"time"
)

func main() {
	var elevAddr string
	var stableID string

	const (
		d            = 500 * time.Millisecond
		brodcastPort = 43452
		TCPPort      = "3000"
	)

	flag.StringVar(&stableID, "id", "", "stable elevator id (A/B/C)")
	flag.StringVar(&elevAddr, "elev", "127.0.0.1:15657", "elevator server addr")
	flag.Parse()

	elevio.Init(elevAddr, 4)
	ps := &types.PeerState{}

	peerID, peerUpdateCh := network.NetworkInit()

	ip, _ := localip.LocalIP()

	e := ElevatorP.NewElevator(peerID, ip)
	e.StableID = stableID

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
				sendmessages.ButtonTransmitLogic(ps, e, btn)
			}

		case newFloor := <-reechFloorCh:
			ElevatorP.ServiceOrderAtFloor(e, newFloor, doorStartTimerCh, ps)

		case <-doorTimeoutCh:
			ElevatorP.OnDoortimeout(doorStartTimerCh, e, ps)

		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %s\n", p.Peers)
			fmt.Printf("  New:      %s\n", p.New)
			fmt.Printf("  Lost:     %s\n", p.Lost)

			RoleManager.RoleElection(p, e, ps, doorStartTimerCh)

			if ps.PrevRole != ps.Role {
				rolechanges.RolesSwitched(ps, TCPPort, TCPRx, e)
				ps.PrevRole = ps.Role

				if ps.Role != types.RolePrimary {
					go tcp.HeartbeatTick(e, ps, 1*time.Second, TCPHeartbeatCh)
					log.Println("entred by stableID:", e.StableID)
				} else {
					go sendmessages.SnapshotTick(e, ps, 200*time.Millisecond)
				}
			}

			rolechanges.HandleLostPeers(p.Lost, e, ps, doorStartTimerCh, p.Peers)

			for stableID, cabs := range types.LostCabOrders {
				log.Println("entred by stableID:", e.StableID)
				log.Printf("lost caborder for stableID: %s, cabs: %v\n", stableID, cabs)
			}

			if len(p.Lost) > 0 && ps.Role == types.RolePrimary && len(p.Peers) > 1 {
				schedueler.MasterSchedueler(e, ps, doorStartTimerCh)
			}

		case PrimaryIdIp := <-UDPHeartbeatRx:
			oldPrimaryID := ps.PrimaryID

			ps.PrimaryID = PrimaryIdIp.PrimaryID
			ps.PrimaryIP = PrimaryIdIp.PrimaryAddrTCP

			if ps.Role != types.RolePrimary &&
				ps.PrimaryID != "" &&
				ps.PrimaryID != oldPrimaryID {
				log.Println("Primary changed, reconnecting to new primary:", ps.PrimaryID)
				go tcp.ConnectToPrimary(ps, TCPPort, e, TCPRx)
			}

		case message := <-TCPRx:
			messagelogic.OnMessageReceive(message, ps, e, doorStartTimerCh)
		}
	}
}
