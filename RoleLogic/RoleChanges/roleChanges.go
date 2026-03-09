package rolechanges

import (
	"heisprosjekt75/Network-go/network/tcp"
	"heisprosjekt75/types"
	
)

func RolesSwitched(ps *types.PeerState, port string, incomingTCP chan tcp.Message, e *types.Elevator) {
	if ps.Role == types.RolePrimary {
		tcp.StartPrimaryTCP(ps, port, incomingTCP,e)
	} else {
		go tcp.ConnectToPrimary(ps, port, e, incomingTCP)
	}
}
