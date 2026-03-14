package messagecomplete

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/Messages/SendMessages"
	"heisprosjekt75/Network-go/network/tcp"
	"heisprosjekt75/types"
	"log"
)

func OrderCompleted(btn elevio.ButtonEvent,e *types.Elevator,ps *types.PeerState){
	messageData := tcp.CompletedOrderMessage{Floor: btn.Floor, Button: btn.Button}
	buttonMessage := tcp.Message{Type: tcp.MsgCompletedOrder, NodeID: e.MyID, MessageData: messageData}
	
	if ps.Role != types.RolePrimary {
		tcp.SendTCP(ps.PrimaryID, buttonMessage, ps)
		log.Printf("Completed order at floor: %d, button: %d, \n", btn.Floor ,btn.Button)
	} else {

		if !types.FullOrderMatrix[btn.Floor][btn.Button] {
		log.Printf("Ignoring duplicate completion for floor:%d button:%d", btn.Floor, btn.Button)
		return
		}

		types.FullOrderMatrix[btn.Floor][btn.Button] = false

		
		for id, matrix := range types.CurrentAssignment {
			matrix[btn.Floor][btn.Button] = false
			types.CurrentAssignment[id] = matrix
		}
		log.Printf("FullOrderMatrix CLEAR -> floor:%d button:%d", btn.Floor, btn.Button)
		sendmessages.SendSnapshotHall(ps,e,types.FullOrderMatrix)
		
	}
}
