package sendmessages

import (
	"fmt"
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/Network-go/network/tcp"
	"heisprosjekt75/types"
)

func SendSnapshot(ps *types.PeerState, e *types.Elevator, hallOrderMAtrix [types.NumFloors][types.NumHallButtons]bool) {
	messageData := tcp.SnapshotHallOrdersMessage{Hall: hallOrderMAtrix}
	buttonMessage := tcp.Message{Type: tcp.MsgSnapshot, NodeID: e.MyID, MessageData: messageData}
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
		fmt.Println("ankommet buttontransmitt logic som ikk master")
	} else {
		types.FullOrderMatrix[btn.Floor][btn.Button] = true
		SendSnapshot(ps, e, types.FullOrderMatrix)

	}
}

func SendHallLightOn(ps *types.PeerState, e *types.Elevator, btn elevio.ButtonEvent, world map[string]types.ElevatorStatus) {
	messageData := tcp.HallLightsOnMessage{Floor: btn.Floor, Button: btn.Button}
	buttonMessage := tcp.Message{Type: tcp.MsgSetHallLights, NodeID: e.MyID, MessageData: messageData}
	for id := range world {
		tcp.SendTCP(id, buttonMessage, ps)
	} //vi skjønner ikke helt om dette vil funke, kan vi se på det sammen på fredag?
}

func SendHallLightOff(ps *types.PeerState, e *types.Elevator, btn elevio.ButtonEvent, world map[string]types.ElevatorStatus) {
	messageData := tcp.HallLightsOffMessage{Floor: btn.Floor, Button: btn.Button}
	buttonMessage := tcp.Message{Type: tcp.MsgTurnOffHallLights, NodeID: e.MyID, MessageData: messageData}
	for id := range world {
		tcp.SendTCP(id, buttonMessage, ps)
	}
}
