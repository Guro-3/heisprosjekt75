package rolechanges

import (
	"heisprosjekt75/Network-go/network/tcp"
	"heisprosjekt75/types"
)

func RolesSwitched(ps *types.PeerState, port string, incomingTCP chan tcp.Message, e *types.Elevator) {
	if ps.Role == types.RolePrimary {
		tcp.StartPrimaryTCP(ps, port, incomingTCP, e)
	} else {
		go tcp.ConnectToPrimary(ps, port, e, incomingTCP)
	}
}

func HandleLostPeers(lost []string, e *types.Elevator, ps *types.PeerState, doorStartTimerCh chan int, currentPeers []string) {
	if ps.Role != types.RolePrimary || len(lost) == 0 {
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
