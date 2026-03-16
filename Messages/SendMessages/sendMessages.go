package sendmessages

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/Messages/MessageTypes"
	"heisprosjekt75/Network-go/network/tcp"
	"heisprosjekt75/types"
)

func SendSnapshot(e *types.Elevator, hallOrderMatrix [types.NumFloors][types.NumHallButtons]bool) {
	if e.Ps.BackupID == "" {
		return
	}
	messageData := messagestypes.SnapshotHallOrdersMessage{Hall: hallOrderMatrix}
	buttonMessage := messagestypes.Message{Type: messagestypes.MsgSnapshot, NodeID: e.MyID, MessageData: messageData}
	tcp.SendTCP(e.Ps.BackupID, buttonMessage, &e.Ps)
}



func SendStateSnapshot(e *types.Elevator) {
	if e.Ps.BackupID == "" {
		return
	}

	worldCopy := make(map[string]types.ElevatorStatus)
	for k, v := range types.WorldView {
		worldCopy[k] = v
	}

	lostCopy := make(map[string][types.NumFloors]bool)
	for k, v := range types.LostCabOrders {
		lostCopy[k] = v
	}

	peerToStableCopy := make(map[string]string)
	for k, v := range types.PeerIDToStableID {
		peerToStableCopy[k] = v
	}

	stableToPeerCopy := make(map[string]string)
	for k, v := range types.StableIDToPeerID {
		stableToPeerCopy[k] = v
	}

	cabCopy := make(map[string][]bool)

	for peerID, state := range types.WorldView {
		cabCopy[peerID] = append([]bool(nil), state.CabRequests...)
	}

	messageData := messagestypes.StateSnapshotMessage{
		Hall:             types.FullOrderMatrix,
		WorldView:        worldCopy,
		LostCabOrders:    lostCopy,
		PeerIDToStableID: peerToStableCopy,
		StableIDToPeerID: stableToPeerCopy,
		CabOrders:        cabCopy,
	}

	buttonMessage := messagestypes.Message{
		Type:        messagestypes.MsgStateSnapshot,
		NodeID:      e.MyID,
		MessageData: messageData,
	}

	tcp.SendTCP(e.Ps.BackupID, buttonMessage, &e.Ps)
}



func SendBackupHallOrderACK(e *types.Elevator) {
	messageData := messagestypes.BackupHallOrderACK{Ack: true}
	buttonMessage := messagestypes.Message{Type: messagestypes.MsgBackupHallOrderACK, NodeID: e.MyID, MessageData: messageData}
	tcp.SendTCP(e.Ps.PrimaryID, buttonMessage, &e.Ps)
}

func ButtonTransmitLogic(e *types.Elevator, btn elevio.ButtonEvent) {
	messageData := messagestypes.HallOrderMessage{Floor: btn.Floor, Button: btn.Button}
	buttonMessage := messagestypes.Message{Type: messagestypes.MsgHallOrder, NodeID: e.MyID, MessageData: messageData}

	if e.Ps.Role != types.RolePrimary {
		tcp.SendTCP(e.Ps.PrimaryID, buttonMessage, &e.Ps)
	} else {
		if !types.FullOrderMatrix[btn.Floor][btn.Button] {
			types.FullOrderMatrix[btn.Floor][btn.Button] = true
			SendStateSnapshot(e)
		}
	}
}



func SendRestoreCabOrders(e *types.Elevator, targetPeerID string, cabs [types.NumFloors]bool) {
	messageData := messagestypes.RestoreCabOrdersMessage{
		NodeID: targetPeerID,
		Cabs:   cabs,
	}

	buttonMessage := messagestypes.Message{
		Type:        messagestypes.MsgRestoreCabOrders,
		NodeID:      e.MyID,
		MessageData: messageData,
	}

	tcp.SendTCP(targetPeerID, buttonMessage, &e.Ps)
}



func SendHallLightOn(e *types.Elevator, btn elevio.ButtonEvent, world map[string]types.ElevatorStatus) {
	messageData := messagestypes.HallLightsOnMessage{Floor: btn.Floor, Button: btn.Button}
	buttonMessage := messagestypes.Message{Type: messagestypes.MsgSetHallLights, NodeID: e.MyID, MessageData: messageData}

	for id := range world {
		if id != e.MyID {
			tcp.SendTCP(id, buttonMessage, &e.Ps)
		}
	}
}



func SendHallLightOff(e *types.Elevator, btn elevio.ButtonEvent, world map[string]types.ElevatorStatus) {
	messageData := messagestypes.HallLightsOffMessage{Floor: btn.Floor, Button: btn.Button}
	buttonMessage := messagestypes.Message{Type: messagestypes.MsgTurnOffHallLights, NodeID: e.MyID, MessageData: messageData}

	for id := range world {
		if id != e.MyID {
			tcp.SendTCP(id, buttonMessage, &e.Ps)
		}
	}
}
