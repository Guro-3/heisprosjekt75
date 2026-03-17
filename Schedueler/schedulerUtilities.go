package schedueler

import (
	"encoding/json"
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/Messages/MessageTypes"
	"heisprosjekt75/Network/tcp"
	"heisprosjekt75/types"
	"log"
	"os/exec"
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

func delegateOrders(receiverID string, e *types.Elevator, btn elevio.ButtonEvent) {
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

func releaseSpecificHallOrder(floor int, button int) {
	for peerID, orderMatrix := range types.CurrentAssignment {
		if orderMatrix[floor][button] {
			orderMatrix[floor][button] = false
			types.CurrentAssignment[peerID] = orderMatrix
		}
	}
}

func getOwnerOfHallOrder(floor int, button int) string {
	for peerID, orderMatrix := range types.CurrentAssignment {
		if orderMatrix[floor][button] {
			return peerID
		}
	}
	return ""
}
