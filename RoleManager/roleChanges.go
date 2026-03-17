package rolemanager

import (
	"heisprosjekt75/Network-go/network/tcp"
	"heisprosjekt75/types"
	"heisprosjekt75/Messages/MessageTypes"
)

func RolesChanges(port string, incomingTCP chan messagestypes.Message, e *types.Elevator) {
	if e.Ps.Role == types.RolePrimary {
		if e.Ps.PrimaryConn != nil {
			e.Ps.PrimaryConn.Close()
		}

		tcp.TcpStartPrimary(port, incomingTCP, e)
	} else {
		go tcp.TcpConnectToPrimary(port, e, incomingTCP)
	}
}


