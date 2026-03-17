package schedueler

import(
	"heisprosjekt75/Network-go/network/tcp"
	"heisprosjekt75/Messages/MessageTypes"
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/types"
	"encoding/json"
	"os/exec"
	"log"
)

func hallAssignmentConvert(matrix [][types.NumHallButtons]bool) types.HallAssignment {
	var out types.HallAssignment

	for f := 0; f < types.NumFloors; f++ {
		for b := 0; b < types.NumHallButtons; b++ {
			if f < len(matrix) {
				out[f][b] = matrix[f][b]
			}
		}
	}
	return out
}



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




func delegateOrders(receiverID string, e *types.Elevator, btn elevio.ButtonEvent, world map[string]types.ElevatorStatus) {
	messageData := messagestypes.HallOrderMessage{Floor: btn.Floor, Button: btn.Button}
	buttonMessage := messagestypes.Message{
		Type:        messagestypes.MsgHallOrder,
		NodeID:      e.MyID,
		MessageData: messageData,
	}
	
	tcp.SendTCP(receiverID, buttonMessage, &e.Ps)
}




func chooseOwner(floor int, button int, proposedAssignment map[string]types.HallAssignment, finalAssignment map[string]types.HallAssignment) string {
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