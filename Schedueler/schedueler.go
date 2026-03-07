package schedueler

import (
	"encoding/json"
	"fmt"
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/Network-go/network/tcp"
	"heisprosjekt75/ElevatorP"
	"heisprosjekt75/types"
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

func DelegateOrders(receiverID string, ps *types.PeerState, e *types.Elevator, btn elevio.ButtonEvent) {
	messageData := tcp.HallOrderMessage{Floor: btn.Floor, Button: btn.Button}
	buttonMessage := tcp.Message{Type: tcp.MsgHallOrder, NodeID: e.MyID, MessageData: messageData}
	tcp.SendTCP(receiverID, buttonMessage, ps)
}

func MasterSchedueler(e *types.Elevator, ps *types.PeerState,doorStartTimerCh chan int) {
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
		fmt.Printf("JSON Marshal error: %v\n", err)
		return
	}
	fmt.Printf("Sending to HRA: %s\n", string(jsonBytes))
	assignment, err := assignHallRequests(jsonBytes)
	if err != nil {
		fmt.Printf("assignHallRequests error: %v\n", err)
		return
	}


	for id, matrix := range assignment {

		var orders []elevio.ButtonEvent

		for f := 0; f < types.NumFloors; f++ {
			for b := 0; b < types.NumHallButtons; b++ {

				if matrix[f][b] {
					orders = append(orders, elevio.ButtonEvent{
						Floor:  f,
						Button: elevio.ButtonType(b),
					})
				}
			}
		}

		types.NodeOrderMap[id] = orders
		fmt.Printf("order to node id %s\n", id)
		if id != e.MyID {
			for _, order := range orders {
				btn := elevio.ButtonEvent{
					Floor:  order.Floor,
					Button: order.Button,
				}
				fmt.Printf("node got order at floor %d and button %d", order.Floor, order.Button)
				DelegateOrders(id, ps, e, btn)
				fmt.Println("Master got completed order")
				types.FullOrderMatrix[ btn.Floor,][btn.Button] = false  // midlertidig plassering

			}
		}
		if id == e.MyID {
			for _, order := range orders {
				btn := elevio.ButtonEvent{
					Floor:  order.Floor,
					Button: order.Button,
				}
				ElevatorP.AddOrder(e, btn.Floor, btn.Button)
				fmt.Printf("I got order at floor %d and button %d", order.Floor, order.Button)
				types.FullOrderMatrix[btn.Floor][ btn.Button] = false  // midlertidig plassering
				ElevatorP.HandleAsignedOrder(e, btn.Floor,  btn.Button, doorStartTimerCh)
			}
			
		}
	}
}
