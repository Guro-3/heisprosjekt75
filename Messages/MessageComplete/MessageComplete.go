package messagecomplete

import (
	"heisprosjekt75/Driver-go/elevio"
	sendmessages "heisprosjekt75/Messages/SendMessages"
	"heisprosjekt75/Network-go/network/tcp"
	"heisprosjekt75/types"
	"log"
)

func OrderCompleted(btn elevio.ButtonEvent, e *types.Elevator, ps *types.PeerState) {
	messageData := tcp.CompletedOrderMessage{Floor: btn.Floor, Button: btn.Button}
	buttonMessage := tcp.Message{Type: tcp.MsgCompletedOrder, NodeID: e.MyID, MessageData: messageData}

	if ps.Role != types.RolePrimary {
		if ps.PrimaryConn == nil {
			log.Printf("FAILED to report completion floor:%d button:%d: no PrimaryConn", btn.Floor, btn.Button)
			return
		}

		tcp.SendTCP(ps.PrimaryID, buttonMessage, ps)
		log.Printf("Completed order at floor:%d button:%d", btn.Floor, btn.Button)
		return
	}

	ApplyCompletedOrder(btn.Floor, btn.Button, e, ps)
}

func ApplyCompletedOrder(floor int, button elevio.ButtonType, e *types.Elevator, ps *types.PeerState) {

	types.FullOrderMatrix[floor][button] = false

	for id, matrix := range types.CurrentAssignment {
		matrix[floor][button] = false
		types.CurrentAssignment[id] = matrix
	}

	log.Printf("FullOrderMatrix CLEAR -> floor:%d button:%d", floor, button)
	sendmessages.SendSnapshot(ps, e, types.FullOrderMatrix)
}
