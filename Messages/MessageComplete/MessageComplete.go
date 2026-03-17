package messagecomplete

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/Messages/MessageTypes"
	"heisprosjekt75/Messages/SendMessages"
	"heisprosjekt75/Network/tcp"
	"heisprosjekt75/types"
	"log"
	"time"
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
		return
	}
	ApplyCompletedOrder(btn.Floor, btn.Button, e)
}

func ApplyCompletedOrder(floor int, button elevio.ButtonType, e *types.Elevator) {
	types.FullOrderMatrix[floor][button] = false
	types.HallOrderTimes[floor][button] = time.Time{}

	for id, matrix := range types.CurrentAssignment {
		matrix[floor][button] = false
		types.CurrentAssignment[id] = matrix
	}
	sendmessages.SendSnapshot(e, types.FullOrderMatrix)
}
