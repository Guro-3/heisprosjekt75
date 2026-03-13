package schedueler

import (
	"encoding/json"
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/ElevatorP"
	"heisprosjekt75/Network-go/network/tcp"
	"heisprosjekt75/types"
	"log"
	"os/exec"
)

func assignHallRequests(input []byte) (map[string][][]bool, error) {

	cmd := exec.Command("./cost_fns/hall_request_assigner/hall_request_assigner", "--input", string(input))

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var result map[string][][]bool

	err = json.Unmarshal(output, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func DelegateOrders(receiverID string, ps *types.PeerState, e *types.Elevator, btn elevio.ButtonEvent, world map[string]types.ElevatorStatus) {
	messageData := tcp.HallOrderMessage{Floor: btn.Floor, Button: btn.Button}
	buttonMessage := tcp.Message{Type: tcp.MsgHallOrder, NodeID: e.MyID, MessageData: messageData}

	tcp.SendTCP(receiverID, buttonMessage, ps)

}

func toHAllAssignment(matrix [][]bool) types.HAllAssignment {
	var out types.HAllAssignment
	for f := 0; f < types.NumFloors; f++ {
		for b := 0; b < types.NumHallButtons; b++ {
			if f < len(matrix) && b < len(matrix[f]) {
				out[f][b] = matrix[f][b]
			}
		}

	}
	return out
}

func MasterSchedueler(e *types.Elevator, ps *types.PeerState, doorStartTimerCh chan int) {
	types.UpdateMyState(e)

	hallRequests := make([][2]bool, types.NumFloors)
	for f := 0; f < types.NumFloors; f++ {
		hallRequests[f] = [2]bool{
			types.FullOrderMatrix[f][0],
			types.FullOrderMatrix[f][1],
		}
	}

	input := types.HRAInput{
		HallRequests: hallRequests,
		States:       types.WorldViewToJSON(types.WorldView),
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		log.Println("JSON Marshal error in MasterSchedueler: ", err)
		return
	}

	assignment, err := assignHallRequests(jsonBytes)
	if err != nil {
		log.Println("assignHallRequests error: ", err)
		return
	}
	newAssignment := make(map[string]types.HAllAssignment)

	for id, matrix := range assignment {
		newAssignment[id] = toHAllAssignment(matrix)
	}

	for id, newMatrix := range newAssignment {
		oldMatrix := types.CurrentAssignment[id]

		for f := 0; f < types.NumFloors; f++ {
			for b := 0; b < types.NumHallButtons; b++ {
				wasAssigned := oldMatrix[f][b]
				isAssigned := newMatrix[f][b]

				if !wasAssigned && isAssigned {
					btn := elevio.ButtonEvent{
						Floor:  f,
						Button: elevio.ButtonType(b),
					}

					if id == e.MyID {
						ElevatorP.HandleAsignedOrder(e, btn.Floor, btn.Button, doorStartTimerCh, ps)

					} else {
						DelegateOrders(id, ps, e, btn, types.WorldView)
					}
				}
			}
		}
	}
	types.CurrentAssignment = newAssignment
}
