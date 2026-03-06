package messagelogic

import (
	"encoding/json"
	"fmt"
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/ElevatorP"
	"heisprosjekt75/Network-go/network/tcp"
	schedueler "heisprosjekt75/Schedueler"
	"heisprosjekt75/types"
	"log"
)

func OnMessageReceive(msg tcp.Message, ps *types.PeerState, e *types.Elevator) {
	switch msg.Type {
	case tcp.MsgHallOrder:
		bytes, _ := json.Marshal(msg.MessageData)

		var order tcp.HallOrderMessage
		json.Unmarshal(bytes, &order)

		log.Printf("Received order at floor: %d, button: %d\n", order.Floor, order.Button)

		switch ps.Role {
		case types.RolePrimary:
			log.Println("Master got hall order")
			types.FullOrderMatrix[order.Floor][order.Button] = true
			schedueler.MasterSchedueler(e, ps)
			fmt.Printf("node got order at floor %d and button %d", order.Floor, order.Button)

		default:
			log.Println("Node got master-message")
			ElevatorP.AddOrder(e, order.Floor, order.Button)
			ElevatorP.StartAction(e)

		}

	case tcp.MsgCompletedOrder:
		bytes, _ := json.Marshal(msg.MessageData)

		var orderComplete tcp.CompletedOrderMessage
		json.Unmarshal(bytes, &orderComplete)

		log.Printf(" completed order at floor: %d, button: %d\n", orderComplete.Floor, orderComplete.Button)

		switch ps.Role {
		case types.RolePrimary:
			log.Println("Master got completed order")
		default:
			log.Println("will not happen")
		}

	case tcp.MsgHeartbeat:
		bytes, _ := json.Marshal(msg.MessageData)

		var heartBeat tcp.HeartbeatMessage
		json.Unmarshal(bytes, &heartBeat)

		log.Printf("Received heartbeat")

		switch ps.Role {
		case types.RolePrimary:
			//log.Printf("Received heartbeat from %s\n", msg.NodeID)
			types.WorldView[msg.NodeID] = types.ElevatorStatus{Floor: heartBeat.CurrentFloor,
				Direction:   heartBeat.Dir,
				State:       heartBeat.State,
				CabRequests: heartBeat.CabRequests}

			//fmt.Printf("state %v\n", heartBeat.State)
		default:
			log.Println("will not happend")
		}

	case tcp.MsgSnapshot:
		bytes, _ := json.Marshal(msg.MessageData)

		var snapshot tcp.SnapshotHallOrdersMessage
		json.Unmarshal(bytes, &snapshot)

		log.Printf("Received snapshot")

		switch ps.Role {
		case types.RoleBackup:
			log.Println("backup received snapshot from master ")
		default:
			log.Println("wil not happend")
		}
	}

}

func ButtonTransmitLogic(ps *types.PeerState, e *types.Elevator, btn elevio.ButtonEvent) {
	messageData := tcp.HallOrderMessage{Floor: btn.Floor, Button: btn.Button}
	buttonMessage := tcp.Message{Type: tcp.MsgHallOrder, NodeID: e.MyID, MessageData: messageData}
	if ps.Role != types.RolePrimary {
		tcp.SendTCP(ps.PrimaryID, buttonMessage, ps)
	} else {
		log.Println("Master got an order!")
	}
}
