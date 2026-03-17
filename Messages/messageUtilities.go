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

func UpdateWorldView(msg messagestypes.Message, worldView messagestypes.WorldViewMessage, e *types.Elevator){
	types.WorldView[msg.NodeID] = types.ElevatorStatus{
				Floor:       worldView.CurrentFloor,
				Direction:   worldView.Dir,
				State:       worldView.State,
				CabRequests: worldView.CabRequests,
			}
			types.TypesUpdateMyState(e)

			if worldView.StableID != "" {
				types.PeerIDToStableID[msg.NodeID] = worldView.StableID
				types.StableIDToPeerID[worldView.StableID] = msg.NodeID

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