package tcp

import(
	"heisprosjekt75/types"
	"heisprosjekt75/Messages/MessageTypes"
)

func handleRestoreCabOrders(e *types.Elevator, peerID string, stableID string) {
	if e.Ps.Role != types.RolePrimary || stableID == "" {
		return
	}

	cabs, ok := types.LostCabOrders[stableID]
	if !ok {
		return
	}

	messageData := messagestypes.RestoreCabOrdersMessage{
		NodeID: peerID,
		Cabs:   cabs,
	}

	buttonMessage := messagestypes.Message{
		Type:        messagestypes.MsgRestoreCabOrders,
		NodeID:      e.MyID,
		MessageData: messageData,
	}

	SendTCP(peerID, buttonMessage, &e.Ps)
	delete(types.LostCabOrders, stableID)
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