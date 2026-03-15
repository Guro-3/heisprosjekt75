package messagelogic

import (
	"encoding/json"
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/ElevatorP"
	sendmessages "heisprosjekt75/Messages/SendMessages"
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

			if !types.FullOrderMatrix[order.Floor][order.Button] {
				types.FullOrderMatrix[order.Floor][order.Button] = true
				sendmessages.SendSnapshot(ps, e, types.FullOrderMatrix)
			}

		default:
			log.Println("got an order at floor:  ", order.Floor, "button: ", order.Button)
			ElevatorP.HandleAsignedOrder(e, order.Floor, order.Button, doorStartTimerCh, ps)

		}

	case tcp.MsgCompletedOrder:
		bytes, _ := json.Marshal(msg.MessageData)

		var orderComplete tcp.CompletedOrderMessage
		json.Unmarshal(bytes, &orderComplete)

		log.Printf("Completed order at floor: %d, button: %d, by elevator: %s\n", orderComplete.Floor, orderComplete.Button, msg.NodeID)

		switch ps.Role {
		case types.RolePrimary:

			if !types.FullOrderMatrix[orderComplete.Floor][orderComplete.Button] {
				log.Printf("Ignoring duplicate completion for floor:%d button:%d",
					orderComplete.Floor, orderComplete.Button)
				return
			}
			types.FullOrderMatrix[orderComplete.Floor][orderComplete.Button] = false

			for id, matrix := range types.CurrentAssignment {
				matrix[orderComplete.Floor][orderComplete.Button] = false
				types.CurrentAssignment[id] = matrix
			}

			log.Printf("FullOrderMatrix CLEAR -> floor:%d button:%d", // ENDRET
				orderComplete.Floor, orderComplete.Button) // ENDRET

			sendmessages.SendSnapshot(ps, e, types.FullOrderMatrix)

		default:
			log.Println("Error: Wrong elevator got MsgCompletedOrder")
		}

	case tcp.MsgHeartbeat:
		bytes, _ := json.Marshal(msg.MessageData)

		var heartBeat tcp.HeartbeatMessage
		json.Unmarshal(bytes, &heartBeat)

		switch ps.Role {
		case types.RolePrimary:
			types.WorldView[msg.NodeID] = types.ElevatorStatus{Floor: heartBeat.CurrentFloor,
				Direction:   heartBeat.Dir,
				State:       heartBeat.State,
				CabRequests: heartBeat.CabRequests}

		default:
			log.Println("Error: Wrong elevator got HeartbeatMessage")
		}

	case tcp.MsgSnapshot:
		bytes, _ := json.Marshal(msg.MessageData)

		var snapshot tcp.SnapshotHallOrdersMessage
		json.Unmarshal(bytes, &snapshot)

		switch ps.Role {
		case types.RoleBackup:
			types.FullOrderMatrix = snapshot.Hall
			sendmessages.BackupHallOrderACK(ps, e)

		default:
			log.Println("Error: Wrong elevator got SnapshotHallOrdersMessage")
		}
	case tcp.MsgBackupHallOrderACK:
		bytes, _ := json.Marshal(msg.MessageData)

		var HallOrderACK tcp.BackupHallOrderACK
		json.Unmarshal(bytes, &HallOrderACK)

		switch ps.Role {
		case types.RolePrimary:
			ElevatorP.SyncHallLight(ps, e, types.WorldView)
			schedueler.MasterSchedueler(e, ps, doorStartTimerCh)

		default:
			log.Println("Error: Wrong elevator got BackupHallOrderACK")
		}

	case tcp.MsgSetHallLights:
		bytes, _ := json.Marshal(msg.MessageData)

		var SetHallLightMsg tcp.HallLightsOnMessage
		json.Unmarshal(bytes, &SetHallLightMsg)

		btn := elevio.ButtonEvent{Floor: SetHallLightMsg.Floor, Button: elevio.ButtonType(SetHallLightMsg.Button)}
		if ps.Role != types.RolePrimary {
			ElevatorP.SetHallLight(btn.Button, btn.Floor)
		}

	case tcp.MsgTurnOffHallLights:
		bytes, _ := json.Marshal(msg.MessageData)

		var TurnOffHallLightMsg tcp.HallLightsOffMessage
		json.Unmarshal(bytes, &TurnOffHallLightMsg)

		btn := elevio.ButtonEvent{Floor: TurnOffHallLightMsg.Floor, Button: elevio.ButtonType(TurnOffHallLightMsg.Button)}
		ElevatorP.TurnOffHallLight(btn.Button, btn.Floor)
	}
}
