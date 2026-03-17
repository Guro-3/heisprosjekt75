package mainutilities

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/Elevator"
	"heisprosjekt75/Messages"
	"heisprosjekt75/Messages/MessageTypes"
	"heisprosjekt75/Messages/SendMessages"
	"heisprosjekt75/Network/SendWorldView"
	"heisprosjekt75/Network/peers"
	"heisprosjekt75/Network/tcp"
	"heisprosjekt75/RoleManager"
	"heisprosjekt75/Schedueler"
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
		schedueler.PrimarySchedueler(e, doorStartTimerCh, -1, -1, "")
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
