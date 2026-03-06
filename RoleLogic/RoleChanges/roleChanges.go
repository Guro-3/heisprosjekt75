package rolechanges

import (
	"heisprosjekt75/types"
	"heisprosjekt75/Network-go/network/tcp"

)

func RolesSwitched(ps *types.PeerState, port string, incomingTCP chan tcp.Message, e *types.Elevator) {
	if ps.Role == types.RolePrimary{
		tcp.StartPrimaryTCP(ps, port, incomingTCP)
	} else {
		tcp.ConnectToPrimary(ps, port, e, incomingTCP)
	}
}