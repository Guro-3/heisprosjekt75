package messagecomplete

import (
	"fmt"
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/Messages/SendMessages"
	"heisprosjekt75/Network-go/network/tcp"
	"heisprosjekt75/types"
)

func OrderCompleted(btn elevio.ButtonEvent,e *types.Elevator,ps *types.PeerState){
	messageData := tcp.CompletedOrderMessage{Floor: btn.Floor, Button: btn.Button}
	buttonMessage := tcp.Message{Type: tcp.MsgCompletedOrder, NodeID: e.MyID, MessageData: messageData}
	if ps.Role != types.RolePrimary {
		tcp.SendTCP(ps.PrimaryID, buttonMessage, ps)
	} else {
		types.FullOrderMatrix[btn.Floor][btn.Button] = false
		sendmessages.SendSnapshot(ps,e,types.FullOrderMatrix)
		fmt.Printf("Master serviced its own order\n")
		
	}
}
