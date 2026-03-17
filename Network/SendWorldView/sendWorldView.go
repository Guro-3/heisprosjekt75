package sendworldview

import (
	"heisprosjekt75/Messages/MessageTypes"
	"heisprosjekt75/Network/tcp"
	"heisprosjekt75/types"
	"time"
)

func WorldViewTick(e *types.Elevator, d time.Duration, TCPWorldViewCh chan<- messagestypes.Message) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for range ticker.C {
		if e.Ps.Role == types.RolePrimary {
			continue
		}

		worldState := messagestypes.WorldViewMessage{
			CurrentFloor: e.CurrentFloor,
			State:        e.State,
			Dir:          e.Dir,
			CabRequests:  e.CabOrderMatrix[:],
			StableID:     e.StableID,
		}

		msg := messagestypes.Message{
			Type:        messagestypes.MsgWorldView,
			NodeID:      e.MyID,
			MessageData: worldState,
		}

		TCPWorldViewCh <- msg
	}
}

func StartWorldViewSender(ps *types.PeerState, WorldViewCh <-chan messagestypes.Message) {
	go func() {
		for msg := range WorldViewCh {
			if ps.PrimaryID != "" {
				tcp.SendTCP(ps.PrimaryID, msg, ps)
			}
		}
	}()
}
