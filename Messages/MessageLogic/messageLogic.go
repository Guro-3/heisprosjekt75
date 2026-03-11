package messagelogic

import (
	"encoding/json"
	"heisprosjekt75/ElevatorP"
	"heisprosjekt75/Messages/SendMessages"
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
			sendmessages.SendSnapshot(ps, e, types.FullOrderMatrix)

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
			sendmessages.SendSnapshot(ps,e,types.FullOrderMatrix)
		default:
			log.Println("shall not happen msg complete order")
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
			log.Println("shall not happen msgHeartbeat")
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
			sendmessages.BackupHallOrderACK(ps, e)

		default:
			log.Println("shall not happend Msgsnapshot")
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
			log.Println("shall not happend MsgBackupHallOrderACK")

		}
	}
}




