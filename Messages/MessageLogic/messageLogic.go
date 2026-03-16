package messagelogic

import (
	"encoding/json"
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/ElevatorP"
	messagecomplete "heisprosjekt75/Messages/MessageComplete"
	sendmessages "heisprosjekt75/Messages/SendMessages"
	"heisprosjekt75/Network-go/network/tcp"
	schedueler "heisprosjekt75/Schedueler"
	"heisprosjekt75/types"
	"log"
)

func OnMessageReceive(msg tcp.Message, ps *types.PeerState, e *types.Elevator, doorStartTimerCh chan int) {
	switch msg.Type {

	case tcp.MsgHallOrder:
		bytes, err := json.Marshal(msg.MessageData)
		if err != nil {
			log.Println("Marshal MsgHallOrder failed:", err)
			return
		}

		var order tcp.HallOrderMessage
		if err := json.Unmarshal(bytes, &order); err != nil {
			log.Println("Unmarshal MsgHallOrder failed:", err)
			return
		}

		switch ps.Role {
		case types.RolePrimary:
			if !types.FullOrderMatrix[order.Floor][order.Button] {
				types.FullOrderMatrix[order.Floor][order.Button] = true
			}

		default:
			log.Printf("Assigned order received -> floor:%d button:%d\n", order.Floor, order.Button)
			ElevatorP.HandleAsignedOrder(e, order.Floor, order.Button, doorStartTimerCh, ps)
		}

	case tcp.MsgCompletedOrder:
		bytes, err := json.Marshal(msg.MessageData)
		if err != nil {
			log.Println("Marshal MsgCompletedOrder failed:", err)
			return
		}

		var orderComplete tcp.CompletedOrderMessage
		if err := json.Unmarshal(bytes, &orderComplete); err != nil {
			log.Println("Unmarshal MsgCompletedOrder failed:", err)
			return
		}

	

		switch ps.Role {
		case types.RolePrimary:
			messagecomplete.ApplyCompletedOrder(orderComplete.Floor, orderComplete.Button, e, ps)

		default:
			log.Println("Error: Wrong elevator got MsgCompletedOrder")
		}

	case tcp.MsgHeartbeat:
		bytes, err := json.Marshal(msg.MessageData)
		if err != nil {
			log.Println("Marshal MsgHeartbeat failed:", err)
			return
		}

		var heartBeat tcp.HeartbeatMessage
		if err := json.Unmarshal(bytes, &heartBeat); err != nil {
			log.Println("Unmarshal MsgHeartbeat failed:", err)
			return
		}

		switch ps.Role {
		case types.RolePrimary:
			types.WorldView[msg.NodeID] = types.ElevatorStatus{
				Floor:       heartBeat.CurrentFloor,
				Direction:   heartBeat.Dir,
				State:       heartBeat.State,
				CabRequests: heartBeat.CabRequests,
			}
			types.UpdateMyState(e)

			if heartBeat.StableID != "" {
				types.PeerIDToStableID[msg.NodeID] = heartBeat.StableID
				types.StableIDToPeerID[heartBeat.StableID] = msg.NodeID

				types.PeerIDToStableID[e.MyID] = e.StableID
				types.StableIDToPeerID[e.StableID] = e.MyID
			}

		default:
			log.Println("Error: Wrong elevator got HeartbeatMessage")
		}

	case tcp.MsgSnapshot:
		bytes, err := json.Marshal(msg.MessageData)
		if err != nil {
			log.Println("Marshal MsgSnapshot failed:", err)
			return
		}

		var snapshot tcp.SnapshotHallOrdersMessage
		if err := json.Unmarshal(bytes, &snapshot); err != nil {
			log.Println("Unmarshal MsgSnapshot failed:", err)
			return
		}

		switch ps.Role {
		case types.RoleBackup:
			types.FullOrderMatrix = snapshot.Hall
			sendmessages.BackupHallOrderACK(ps, e)

		default:
			log.Println("Error: Wrong elevator got SnapshotHallOrdersMessage")
		}
	case tcp.MsgStateSnapshot:
		bytes, err := json.Marshal(msg.MessageData)
		if err != nil {
			log.Println("Marshal MsgStateSnapshot failed:", err)
			return
		}

		var snapshot tcp.StateSnapshotMessage
		if err := json.Unmarshal(bytes, &snapshot); err != nil {
			log.Println("Unmarshal MsgStateSnapshot failed:", err)
			return
		}

		switch ps.Role {
		case types.RoleBackup:
			types.FullOrderMatrix = snapshot.Hall
			types.WorldView = snapshot.WorldView
			types.LostCabOrders = snapshot.LostCabOrders
			types.PeerIDToStableID = snapshot.PeerIDToStableID
			types.StableIDToPeerID = snapshot.StableIDToPeerID

			for peerID, cabs := range snapshot.CabOrders {
				state := types.WorldView[peerID]
				state.CabRequests = cabs
				types.WorldView[peerID] = state
			}

			sendmessages.BackupHallOrderACK(ps, e)

		default:
			log.Println("Error: Wrong elevator got MsgStateSnapshot")
		}

	case tcp.MsgRestoreCabOrders:
		bytes, err := json.Marshal(msg.MessageData)
		if err != nil {
			log.Println("Marshal MsgRestoreCabOrders failed:", err)
			return
		}

		var restore tcp.RestoreCabOrdersMessage
		if err := json.Unmarshal(bytes, &restore); err != nil {
			log.Println("Unmarshal MsgRestoreCabOrders failed:", err)
			return
		}

		if restore.NodeID != e.MyID {
			return
		}

		for f := 0; f < types.NumFloors; f++ {
			if restore.Cabs[f] {
				e.CabOrderMatrix[f] = true
				ElevatorP.SetCabLight(f)
			}
		}

		ElevatorP.StartAction(e, doorStartTimerCh, ps)

	case tcp.MsgBackupHallOrderACK:
		bytes, err := json.Marshal(msg.MessageData)
		if err != nil {
			log.Println("Marshal MsgBackupHallOrderACK failed:", err)
			return
		}

		var hallOrderACK tcp.BackupHallOrderACK
		if err := json.Unmarshal(bytes, &hallOrderACK); err != nil {
			log.Println("Unmarshal MsgBackupHallOrderACK failed:", err)
			return
		}

		switch ps.Role {
		case types.RolePrimary:
			if hallOrderACK.Ack {
				ElevatorP.SyncHallLight(ps, e, types.WorldView)
				schedueler.MasterSchedueler(e, ps, doorStartTimerCh)
			}

		default:
			log.Println("Error: Wrong elevator got BackupHallOrderACK")
		}

	case tcp.MsgSetHallLights:
		bytes, err := json.Marshal(msg.MessageData)
		if err != nil {
			log.Println("Marshal MsgSetHallLights failed:", err)
			return
		}

		var setHallLightMsg tcp.HallLightsOnMessage
		if err := json.Unmarshal(bytes, &setHallLightMsg); err != nil {
			log.Println("Unmarshal MsgSetHallLights failed:", err)
			return
		}

		btn := elevio.ButtonEvent{
			Floor:  setHallLightMsg.Floor,
			Button: elevio.ButtonType(setHallLightMsg.Button),
		}

		if ps.Role != types.RolePrimary {
			ElevatorP.SetHallLight(btn.Button, btn.Floor)
		}

	case tcp.MsgTurnOffHallLights:
		bytes, err := json.Marshal(msg.MessageData)
		if err != nil {
			log.Println("Marshal MsgTurnOffHallLights failed:", err)
			return
		}

		var turnOffHallLightMsg tcp.HallLightsOffMessage
		if err := json.Unmarshal(bytes, &turnOffHallLightMsg); err != nil {
			log.Println("Unmarshal MsgTurnOffHallLights failed:", err)
			return
		}

		btn := elevio.ButtonEvent{
			Floor:  turnOffHallLightMsg.Floor,
			Button: elevio.ButtonType(turnOffHallLightMsg.Button),
		}

		ElevatorP.TurnOffHallLight(btn.Button, btn.Floor)
	}
}
