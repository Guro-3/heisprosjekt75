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

func assignHallRequests(input []byte) (map[string][][types.NumHallButtons]bool, error) {

	cmd := exec.Command("./cost_fns/hall_request_assigner/hall_request_assigner", "-i", string(input))

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("hall_request_assigner output:", string(output))
		return nil, err
	}

	var result map[string][][types.NumHallButtons]bool

	err = json.Unmarshal(output, &result)
	if err != nil {
		log.Println("hall_request_assigner raw output:", string(output))
		return nil, err
	}

	return result, nil
}

func DelegateOrders(receiverID string, ps *types.PeerState, e *types.Elevator, btn elevio.ButtonEvent, world map[string]types.ElevatorStatus) {

	messageData := tcp.HallOrderMessage{Floor: btn.Floor, Button: btn.Button}

	buttonMessage := tcp.Message{
		Type:        tcp.MsgHallOrder,
		NodeID:      e.MyID,
		MessageData: messageData,
	}

	tcp.SendTCP(receiverID, buttonMessage, ps)
}

func toHAllAssignment(matrix [][2]bool) types.HAllAssignment {

	var out types.HAllAssignment

	for f := 0; f < types.NumFloors; f++ {
		for b := 0; b < types.NumHallButtons; b++ {
			if f < len(matrix) {
				out[f][b] = matrix[f][b]
			}
		}
	}

	return out

}

func chooseOwner(floor int, button int, proposedAssignment map[string]types.HAllAssignment, finalAssignment map[string]types.HAllAssignment) string {

	owner := ""
	other := 1 - button

	for id, matrix := range types.CurrentAssignment {
		_, alive := types.WorldView[id]
		if alive && matrix[floor][button] && !finalAssignment[id][floor][other] {
			owner = id
			break
		}
	}

	if owner == "" {
		for id, matrix := range proposedAssignment {

			_, alive := types.WorldView[id]
			if alive && matrix[floor][button] && !finalAssignment[id][floor][other] {
				owner = id
				break
			}
		}
	}

	if owner == "" {
		for id := range types.WorldView {
			_, alive := types.WorldView[id]
			if alive && !finalAssignment[id][floor][other] {
				owner = id
				break
			}
		}
	}

	return owner
}

func MasterSchedueler(e *types.Elevator, ps *types.PeerState, doorStartTimerCh chan int) {

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
		log.Println("JSON Marshal error in MasterSchedueler:", err)
		return
	}

	proposal, err := assignHallRequests(jsonBytes)
	if err != nil {
		log.Println("assignHallRequests error:", err)
		return
	}

	proposedAssignment := make(map[string]types.HAllAssignment)

	for id, matrix := range proposal {
		proposedAssignment[id] = toHAllAssignment(matrix)
	}

	finalAssignment := make(map[string]types.HAllAssignment)

	for id := range types.WorldView {
		finalAssignment[id] = types.HAllAssignment{}
	}

	for f := 0; f < types.NumFloors; f++ {
		for b := 0; b < types.NumHallButtons; b++ {

			if !types.FullOrderMatrix[f][b] {
				continue
			}

			owner := chooseOwner(f, b, proposedAssignment, finalAssignment)

			if owner != "" {
				m := finalAssignment[owner]
				m[f][b] = true
				finalAssignment[owner] = m
			}
		}
	}

	for id, newMatrix := range finalAssignment {

		oldMatrix := types.CurrentAssignment[id]

		for f := 0; f < types.NumFloors; f++ {
			for b := 0; b < types.NumHallButtons; b++ {

				if !oldMatrix[f][b] && newMatrix[f][b] {

					btn := elevio.ButtonEvent{
						Floor:  f,
						Button: elevio.ButtonType(b),
					}

					log.Printf("Assign -> %s floor:%d button:%d", id, f, b)

					if id == e.MyID {
						ElevatorP.HandleAsignedOrder(e, f, elevio.ButtonType(b), doorStartTimerCh, ps)
					} else {
						DelegateOrders(id, ps, e, btn, types.WorldView)
					}
				}
			}
		}
	}

	types.CurrentAssignment = finalAssignment
}
