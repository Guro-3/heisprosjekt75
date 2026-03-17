package messagelogic

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/Elevator"
	"heisprosjekt75/Messages"
	"heisprosjekt75/Messages/MessageComplete"
	"heisprosjekt75/Messages/MessageTypes"
	"heisprosjekt75/Messages/SendMessages"
	"heisprosjekt75/Schedueler"
	"heisprosjekt75/types"
	"log"
)

func OnMessageReceive(msg messagestypes.Message, e *types.Elevator, doorStartTimerCh chan int) {
	switch msg.Type {
	case messagestypes.MsgHallOrder:
		order, err := messages.DecodeMessage[messagestypes.HallOrderMessage](msg.MessageData)
		if err != nil {
			log.Println("Decode MsgHallOrder failed:", err)
			return
		}
		switch e.Ps.Role {
		case types.RolePrimary:
			if !types.FullOrderMatrix[order.Floor][order.Button] {
				types.FullOrderMatrix[order.Floor][order.Button] = true
			}

		default:
			Elevator.HandleAssignedOrder(e, order.Floor, order.Button, doorStartTimerCh)
		}

	case messagestypes.MsgCompletedOrder:
		orderComplete, err := messages.DecodeMessage[messagestypes.CompletedOrderMessage](msg.MessageData)
		if err != nil {
			log.Println("Decode MsgCompletedOrder failed:", err)
			return
		}

		switch e.Ps.Role {
		case types.RolePrimary:
			messagecomplete.ApplyCompletedOrder(orderComplete.Floor, orderComplete.Button, e)

		default:
			log.Println("Error: Wrong elevator got MsgCompletedOrder")
		}

	case messagestypes.MsgWorldView:
		worldView, err := messages.DecodeMessage[messagestypes.WorldViewMessage](msg.MessageData)
		if err != nil {
			log.Println("Decode MsgWorldView failed:", err)
			return
		}

		switch e.Ps.Role {
		case types.RolePrimary:
			messages.UpdateWorldView(msg, worldView, e)
		default:
			log.Println("Error: Wrong elevator got WorldViewMessage")
		}

	case messagestypes.MsgSnapshot:
		snapshot, err := messages.DecodeMessage[messagestypes.SnapshotHallOrdersMessage](msg.MessageData)
		if err != nil {
			log.Println("Decode MsgSnapshot failed:", err)
			return
		}

		switch e.Ps.Role {
		case types.RoleBackup:
			types.FullOrderMatrix = snapshot.Hall
			sendmessages.SendBackupHallOrderACK(e)

		default:
			log.Println("Error: Wrong elevator got SnapshotHallOrdersMessage")
		}

	case messagestypes.MsgStateSnapshot:
		snapshot, err := messages.DecodeMessage[messagestypes.StateSnapshotMessage](msg.MessageData)
		if err != nil {
			log.Println("Decode MsgStateSnapshot failed:", err)
			return
		}

		switch e.Ps.Role {
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

			sendmessages.SendBackupHallOrderACK(e)

		default:
			log.Println("Error: Wrong elevator got MsgStateSnapshot")
		}

	case messagestypes.MsgRestoreCabOrders:
		restore, err := messages.DecodeMessage[messagestypes.RestoreCabOrdersMessage](msg.MessageData)
		if err != nil {
			log.Println("Decode MsgRestoreCabOrders failed:", err)
			return
		}

		if restore.NodeID != e.MyID {
			return
		}

		for f := 0; f < types.NumFloors; f++ {
			if restore.Cabs[f] {
				e.CabOrderMatrix[f] = true
				Elevator.SetCabLight(f)
			}
		}

		Elevator.FsmStartAction(e, doorStartTimerCh)

	case messagestypes.MsgBackupHallOrderACK:
		hallOrderACK, err := messages.DecodeMessage[messagestypes.BackupHallOrderACK](msg.MessageData)
		if err != nil {
			log.Println("Decode MsgBackupHallOrderACKfailed:", err)
			return
		}

		switch e.Ps.Role {
		case types.RolePrimary:
			if hallOrderACK.Ack {
				Elevator.SyncHallLight(e, types.WorldView)
				schedueler.PrimarySchedueler(e, doorStartTimerCh)
			}

		default:
			log.Println("Error: Wrong elevator got BackupHallOrderACK")
		}

	case messagestypes.MsgSetHallLights:
		setHallLightMsg, err := messages.DecodeMessage[messagestypes.HallLightsOnMessage](msg.MessageData)
		if err != nil {
			log.Println("Decode MsgSetHallLights:", err)
			return
		}

		btn := elevio.ButtonEvent{
			Floor:  setHallLightMsg.Floor,
			Button: elevio.ButtonType(setHallLightMsg.Button),
		}

		if e.Ps.Role != types.RolePrimary {
			Elevator.SetHallLight(btn.Button, btn.Floor)
		}

	case messagestypes.MsgTurnOffHallLights:
		turnOffHallLightMsg, err := messages.DecodeMessage[messagestypes.HallLightsOffMessage](msg.MessageData)
		if err != nil {
			log.Println("Decode MsgTurnOffHallLights:", err)
			return
		}


		btn := elevio.ButtonEvent{
			Floor:  turnOffHallLightMsg.Floor,
			Button: elevio.ButtonType(turnOffHallLightMsg.Button),
		}

		Elevator.TurnOffHallLight(btn.Button, btn.Floor)
	}
}

