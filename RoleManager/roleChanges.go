package rolemanager

import (
	messagestypes "heisprosjekt75/Messages/MessageTypes"
	"heisprosjekt75/Network/tcp"
	"heisprosjekt75/types"
)

func RolesChanges(port string, incomingTCP chan messagestypes.Message, e *types.Elevator) {
	/*if e.Ps.Role != types.RolePrimary {
		if e.Ps.PrimaryListener != nil {
			e.Ps.PrimaryListener.Close()
			e.Ps.PrimaryListener = nil
		}
	}

	if e.Ps.PrimaryConn != nil {
		e.Ps.PrimaryConn.Close()
		e.Ps.PrimaryConn = nil
	}*/

	if e.Ps.Role == types.RolePrimary {
		tcp.TcpStartPrimary(port, incomingTCP, e)
	} else {
		go tcp.TcpConnectToPrimary(port, e, incomingTCP)
	}
}
