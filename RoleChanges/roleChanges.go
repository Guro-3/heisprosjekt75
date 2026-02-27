package rolechanges

import (
	"heisprosjekt75/RoleManager"
	"heisprosjekt75/ElevatorP"
	"heisprosjekt75/Network-go/network/tcp"

)

func RolesSwitched(ps *RoleManager.PeerState, port string, incomingTCP chan string, e *ElevatorP.Elevator) {
	if ps.Role == RoleManager.RolePrimary{
		tcp.StartPrimaryTCP(ps, port, incomingTCP)
	} else {
		tcp.ConnectToPrimary(ps, port, e, incomingTCP)
	}
}