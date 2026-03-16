package messagecomplete

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/Messages/MessageTypes"
	"heisprosjekt75/Messages/SendMessages"
	"heisprosjekt75/Network-go/network/tcp"
	"heisprosjekt75/types"
	"log"
)

func OrderCompleted(btn elevio.ButtonEvent, e *types.Elevator) {
	messageData := messagestypes.CompletedOrderMessage{Floor: btn.Floor, Button: btn.Button}
	buttonMessage := messagestypes.Message{Type: messagestypes.MsgCompletedOrder, NodeID: e.MyID, MessageData: messageData}

	if e.Ps.Role != types.RolePrimary {
		if e.Ps.PrimaryConn == nil {
			log.Printf("FAILED to report completion floor:%d button:%d: no PrimaryConn", btn.Floor, btn.Button)
			return
		}
		tcp.SendTCP(e.Ps.PrimaryID, buttonMessage, &e.Ps)
		log.Printf("Completed order at floor:%d button:%d", btn.Floor, btn.Button)
		return
	}
	ApplyCompletedOrder(btn.Floor, btn.Button, e)
}



func ApplyCompletedOrder(floor int, button elevio.ButtonType, e *types.Elevator) {
	types.FullOrderMatrix[floor][button] = false

	for id, matrix := range types.CurrentAssignment {
		matrix[floor][button] = false
		types.CurrentAssignment[id] = matrix
	}

	log.Printf("FullOrderMatrix CLEAR -> floor:%d button:%d", floor, button)
	sendmessages.SendSnapshot(e, types.FullOrderMatrix)
}
