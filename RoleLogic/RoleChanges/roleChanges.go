package rolechanges

import (
	"heisprosjekt75/Network-go/network/tcp"
	"heisprosjekt75/types"
	"heisprosjekt75/Messages/MessageTypes"
)

func RolesSwitched(port string, incomingTCP chan messagestypes.Message, e *types.Elevator) {
	if e.Ps.Role == types.RolePrimary {
		if e.Ps.PrimaryConn != nil {
			e.Ps.PrimaryConn.Close()
		}

		tcp.StartPrimaryTCP(port, incomingTCP, e)
	} else {
		go tcp.ConnectToPrimary(port, e, incomingTCP)
	}
}

func HandleLostPeers(lost []string, e *types.Elevator, doorStartTimerCh chan int, currentPeers []string) {
	if e.Ps.Role != types.RolePrimary || len(lost) == 0 {
		return
	}

	for _, lostPeerID := range lost {
		if st, ok := types.WorldView[lostPeerID]; ok {
			if stableID, ok := types.PeerIDToStableID[lostPeerID]; ok && stableID != "" {
				var cabCopy [types.NumFloors]bool
				for i := 0; i < types.NumFloors && i < len(st.CabRequests); i++ {
					cabCopy[i] = st.CabRequests[i]
				}
				types.LostCabOrders[stableID] = cabCopy
			}
		}

		if stableID, ok := types.PeerIDToStableID[lostPeerID]; ok {
			delete(types.StableIDToPeerID, stableID)
		}

		delete(types.PeerIDToStableID, lostPeerID)
		delete(types.WorldView, lostPeerID)
		delete(types.CurrentAssignment, lostPeerID)
	}

}
