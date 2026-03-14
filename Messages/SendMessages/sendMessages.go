package sendmessages

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/Network-go/network/tcp"
	"heisprosjekt75/types"
)

func SendSnapshotHall(ps *types.PeerState, e *types.Elevator, hallOrderMAtrix [types.NumFloors][types.NumHallButtons]bool) {
	messageData := tcp.SnapshotHallOrdersMessage{Hall: hallOrderMAtrix}
	buttonMessage := tcp.Message{Type: tcp.MsgSnapshotHall, NodeID: e.MyID, MessageData: messageData}
	tcp.SendTCP(ps.BackupID, buttonMessage, ps)
}

func SendSnapshotCabs(ps *types.PeerState, e *types.Elevator, CabOrderMatrix map[string]types.CabOrderMatrix) {
	messageData := tcp.SnapshotCabOrdersMessage{Cabs: CabOrderMatrix}
	buttonMessage := tcp.Message{Type: tcp.MsgSnapshotCabs, NodeID: e.MyID, MessageData: messageData}
	tcp.SendTCP(ps.BackupID, buttonMessage, ps)
}

func SendCabOrdersToPrimary(ps *types.PeerState, e *types.Elevator, ActiveCabOrders [types.NumFloors]bool) {
	messageData := tcp.CabOrderMessage{Cabs: ActiveCabOrders, NodeIP: e.ElevIP}
	buttonMessage := tcp.Message{Type: tcp.MsgCabOrders, NodeID: e.MyID, MessageData: messageData}
	tcp.SendTCP(ps.PrimaryID, buttonMessage, ps)
}

func SendCabOrdersToNode(ps *types.PeerState, e *types.Elevator, ActiveCabOrders [types.NumFloors]bool, receiverID string) {
	messageData := tcp.CabOrderMessage{Cabs: ActiveCabOrders, NodeIP: e.ElevIP}
	buttonMessage := tcp.Message{Type: tcp.MsgCabOrders, NodeID: e.MyID, MessageData: messageData}
	tcp.SendTCP(receiverID, buttonMessage, ps)
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
			SendSnapshotHall(ps, e, types.FullOrderMatrix)
		}
	}
}

func SendHallLightOn(ps *types.PeerState, e *types.Elevator, btn elevio.ButtonEvent, world map[string]types.ElevatorStatus) {
	messageData := tcp.HallLightsOnMessage{Floor: btn.Floor, Button: btn.Button}
	buttonMessage := tcp.Message{Type: tcp.MsgSetHallLights, NodeID: e.MyID, MessageData: messageData}
	for id:= range world {
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
