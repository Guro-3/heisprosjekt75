package messages

import(
	"encoding/json"
	"heisprosjekt75/types"
	"heisprosjekt75/Messages/MessageTypes"
	"time"
	"heisprosjekt75/Messages/SendMessages"
)

func DecodeMessage[T any](data interface{}) (T, error) {
	var result T

	bytes, err := json.Marshal(data)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(bytes, &result)
	return result, err
}

func UpdateWorldView(msg messagestypes.Message, heartBeat messagestypes.HeartbeatMessage, e *types.Elevator){
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
}

func SnapshotTick(e *types.Elevator, d time.Duration) {
	tic := time.NewTicker(d)
	defer tic.Stop()

	for range tic.C {
		if e.Ps.Role != types.RolePrimary {
			continue
		}

		sendmessages.SendStateSnapshot(e)
	}
}