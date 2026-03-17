package mainutilities

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/Elevator"
	messages "heisprosjekt75/Messages"
	messagestypes "heisprosjekt75/Messages/MessageTypes"
	sendmessages "heisprosjekt75/Messages/SendMessages"
	sendworldview "heisprosjekt75/Network-go/network/SendWorldView"
	"heisprosjekt75/Network-go/network/peers"
	"heisprosjekt75/Network-go/network/tcp"
	rolemanager "heisprosjekt75/RoleManager"
	schedueler "heisprosjekt75/Schedueler"
	"heisprosjekt75/types"
	"time"
)

func CaseReceiveBtnCh(btn elevio.ButtonEvent, e *types.Elevator, doorStartTimerCh chan int) {
	if e.Mode == types.SingleElevator || btn.Button == elevio.BT_Cab {
		Elevator.FsmServiceLocalButton(e, btn.Floor, btn.Button, doorStartTimerCh)
	} else {
		sendmessages.ButtonTransmitLogic(e, btn)
	}
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func CasePeerUpdateCh(e *types.Elevator, doorStartTimerCh chan int, TCPPort string, TCPRx chan messagestypes.Message, p peers.PeerUpdate, TCPWorldViewCh chan<- messagestypes.Message) {

	rolemanager.RoleElection(p, e, doorStartTimerCh)

	if e.Ps.PrevRole != e.Ps.Role {
		rolemanager.RolesChanges(TCPPort, TCPRx, e)
		e.Ps.PrevRole = e.Ps.Role

		if e.Ps.Role != types.RolePrimary {
			go sendworldview.WorldViewTick(e, 1*time.Second, TCPWorldViewCh)

		} else {
			go messages.SnapshotTick(e, 500*time.Millisecond)
		}
	}

	tcp.HandleLostPeers(p.Lost, e, doorStartTimerCh, p.Peers)

	if len(p.Lost) > 0 && e.Ps.Role == types.RolePrimary && len(p.Peers) > 1 {
		schedueler.PrimarySchedueler(e, doorStartTimerCh)
	}
}

func CaseUDPPrimaryIPIDRx(e *types.Elevator, TCPPort string, TCPRx chan messagestypes.Message, PrimaryIdIp sendworldview.PrimaryIPID) {

	oldPrimaryID := e.Ps.PrimaryID
	e.Ps.PrimaryID = PrimaryIdIp.PrimaryID
	e.Ps.PrimaryIP = PrimaryIdIp.PrimaryAddrTCP

	if e.Ps.Role != types.RolePrimary &&

		e.Ps.PrimaryID != "" &&
		e.Ps.PrimaryID != oldPrimaryID {
		go tcp.TcpConnectToPrimary(TCPPort, e, TCPRx)
	}

}
