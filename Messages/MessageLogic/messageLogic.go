package messagelogic

import (
	"encoding/json"
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/ElevatorP"
	"heisprosjekt75/Network-go/network/tcp"
	schedueler "heisprosjekt75/Schedueler"
	"heisprosjekt75/types"
	"log"
)

func OnMessageReceive(msg tcp.Message, ps *types.PeerState, e *types.Elevator, doorStartTimerCh chan int) {
	switch msg.Type {
	case tcp.MsgHallOrder:
		bytes, _ := json.Marshal(msg.MessageData)

		var order tcp.HallOrderMessage
		json.Unmarshal(bytes, &order)

		switch ps.Role {
		case types.RolePrimary:

			types.FullOrderMatrix[order.Floor][order.Button] = true

		default:

			ElevatorP.HandleAsignedOrder(e, order.Floor, order.Button, doorStartTimerCh, ps)

		}

	case tcp.MsgCompletedOrder:
		bytes, _ := json.Marshal(msg.MessageData)

		var orderComplete tcp.CompletedOrderMessage
		json.Unmarshal(bytes, &orderComplete)

		log.Printf("Completed order at floor: %d, button: %d\n", orderComplete.Floor, orderComplete.Button)

		switch ps.Role {
		case types.RolePrimary:
			log.Println("Master got completed order")
			types.FullOrderMatrix[orderComplete.Floor][orderComplete.Button] = false
		default:
			log.Println("will not happen")
		}

	case tcp.MsgHeartbeat:
		bytes, _ := json.Marshal(msg.MessageData)

		var heartBeat tcp.HeartbeatMessage
		json.Unmarshal(bytes, &heartBeat)

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
			log.Println("backup received snapshot from master")
			types.FullOrderMatrix = snapshot.Hall
			BackupHallOrderACK(ps, e)

		default:
			log.Println("wil not happend")
		}
	case tcp.MsgBackupHallOrderACK:
		bytes, _ := json.Marshal(msg.MessageData)

		var HallOrderACK tcp.BackupHallOrderACK
		json.Unmarshal(bytes, &HallOrderACK)

		log.Printf("Received snapshot")

		switch ps.Role {
		case types.RolePrimary:
			schedueler.MasterSchedueler(e, ps, doorStartTimerCh)
		default:
			log.Println("wil not happend")

		}
	}
}

func ButtonTransmitLogic(ps *types.PeerState, e *types.Elevator, btn elevio.ButtonEvent, doorStartTimerCh chan int) {
	messageData := tcp.HallOrderMessage{Floor: btn.Floor, Button: btn.Button}
	buttonMessage := tcp.Message{Type: tcp.MsgHallOrder, NodeID: e.MyID, MessageData: messageData}
	if ps.Role != types.RolePrimary {
		tcp.SendTCP(ps.PrimaryID, buttonMessage, ps)
	} else {
		types.FullOrderMatrix[btn.Floor][btn.Button] = true
		schedueler.MasterSchedueler(e, ps, doorStartTimerCh)
	}
}

func BackupHallOrderACK(ps *types.PeerState, e *types.Elevator) {
	messageData := tcp.BackupHallOrderACK{Ack: true}
	buttonMessage := tcp.Message{Type: tcp.MsgBackupHallOrderACK, NodeID: e.MyID, MessageData: messageData}
	tcp.SendTCP(ps.PrimaryID, buttonMessage, ps)
}
