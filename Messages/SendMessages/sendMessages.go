package sendmessages

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/Network-go/network/tcp"
	"heisprosjekt75/types"
	"log"
)

func SendSnapshot(ps *types.PeerState, e *types.Elevator, hallOrderMatrix [types.NumFloors][types.NumHallButtons]bool) {
	if ps.BackupID == "" {
		return
	}

	messageData := tcp.SnapshotHallOrdersMessage{Hall: hallOrderMatrix}
	buttonMessage := tcp.Message{Type: tcp.MsgSnapshot, NodeID: e.MyID, MessageData: messageData}
	tcp.SendTCP(ps.BackupID, buttonMessage, ps)
}

func SendStateSnapshot(ps *types.PeerState, e *types.Elevator) {
	if ps.BackupID == "" {
		return
	}
	log.Println("entered sendStateStapshot som id:", e.StableID)

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

	messageData := tcp.StateSnapshotMessage{
		Hall:             types.FullOrderMatrix,
		WorldView:        worldCopy,
		LostCabOrders:    lostCopy,
		PeerIDToStableID: peerToStableCopy,
		StableIDToPeerID: stableToPeerCopy,
	}

	buttonMessage := tcp.Message{
		Type:        tcp.MsgStateSnapshot,
		NodeID:      e.MyID,
		MessageData: messageData,
	}

	tcp.SendTCP(ps.BackupID, buttonMessage, ps)
}

func BackupHallOrderACK(ps *types.PeerState, e *types.Elevator) {
	messageData := tcp.BackupHallOrderACK{Ack: true}
	buttonMessage := tcp.Message{Type: tcp.MsgBackupHallOrderACK, NodeID: e.MyID, MessageData: messageData}
	tcp.SendTCP(ps.PrimaryID, buttonMessage, ps)
}

func ButtonTransmitLogic(ps *types.PeerState, e *types.Elevator, btn elevio.ButtonEvent) {
	messageData := tcp.HallOrderMessage{Floor: btn.Floor, Button: btn.Button}
	buttonMessage := tcp.Message{Type: tcp.MsgHallOrder, NodeID: e.MyID, MessageData: messageData}

	if ps.Role != types.RolePrimary {
		tcp.SendTCP(ps.PrimaryID, buttonMessage, ps)
	} else {
		if !types.FullOrderMatrix[btn.Floor][btn.Button] {
			types.FullOrderMatrix[btn.Floor][btn.Button] = true
			SendStateSnapshot(ps, e)
		}
	}
}

func SendRestoreCabOrders(ps *types.PeerState, e *types.Elevator, targetPeerID string, cabs [types.NumFloors]bool) {
	messageData := tcp.RestoreCabOrdersMessage{
		NodeID: targetPeerID,
		Cabs:   cabs,
	}

	buttonMessage := tcp.Message{
		Type:        tcp.MsgRestoreCabOrders,
		NodeID:      e.MyID,
		MessageData: messageData,
	}
	log.Printf("er inni sendRestoreCaborders")
	tcp.SendTCP(targetPeerID, buttonMessage, ps)
}

func SendHallLightOn(ps *types.PeerState, e *types.Elevator, btn elevio.ButtonEvent, world map[string]types.ElevatorStatus) {
	messageData := tcp.HallLightsOnMessage{Floor: btn.Floor, Button: btn.Button}
	buttonMessage := tcp.Message{Type: tcp.MsgSetHallLights, NodeID: e.MyID, MessageData: messageData}

	for id := range world {
		if id != e.MyID {
			tcp.SendTCP(id, buttonMessage, ps)
		}
	}
}

func SendHallLightOff(ps *types.PeerState, e *types.Elevator, btn elevio.ButtonEvent, world map[string]types.ElevatorStatus) {
	messageData := tcp.HallLightsOffMessage{Floor: btn.Floor, Button: btn.Button}
	buttonMessage := tcp.Message{Type: tcp.MsgTurnOffHallLights, NodeID: e.MyID, MessageData: messageData}

	for id := range world {
		if id != e.MyID {
			tcp.SendTCP(id, buttonMessage, ps)
		}
	}
}
