package schedueler

import (
	"encoding/json"
	"fmt"
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/ElevatorP"
	sendmessages "heisprosjekt75/Messages/SendMessages"
	"heisprosjekt75/Network-go/network/tcp"
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

func DelegateOrders(receiverID string, ps *types.PeerState, e *types.Elevator, btn elevio.ButtonEvent, world map[string]types.ElevatorStatus) {
	messageData := tcp.HallOrderMessage{Floor: btn.Floor, Button: btn.Button}
	buttonMessage := tcp.Message{Type: tcp.MsgHallOrder, NodeID: e.MyID, MessageData: messageData}
	fmt.Printf("Deligate orders\n")
	tcp.SendTCP(receiverID, buttonMessage, ps)
	sendmessages.SendHallLightOn(ps, e, btn, world)
	ElevatorP.SetHallLight(btn.Button, btn.Floor)
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
	fmt.Printf("hallrquest: %v\n", hallRequests)

	input := types.HRAInput{
		HallRequests: hallRequests,
		States:       types.WorldViewToJSON(types.WorldView),
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Printf("JSON Marshal error: %v\n", err)
		return
	}

	assignment, err := assignHallRequests(jsonBytes)
	if err != nil {
		fmt.Printf("assignHallRequests error: %v\n", err)
		return
	}

	for id, matrix := range assignment {
		fmt.Printf("reciever id: %s\n", id)
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

		fmt.Printf("order to node id %s\n", id)
		if id != e.MyID {
			for _, order := range orders {
				btn := elevio.ButtonEvent{
					Floor:  order.Floor,
					Button: order.Button,
				}

				DelegateOrders(id, ps, e, btn, types.WorldView)

			}
		}
		if id == e.MyID {
			fmt.Printf("Schedulere: this is my id: %s\n", id)
			for _, order := range orders {
				btn := elevio.ButtonEvent{
					Floor:  order.Floor,
					Button: order.Button,
				}
				sendmessages.SendHallLightOn(ps, e, btn, types.WorldView)
				ElevatorP.SetHallLight(btn.Button, btn.Floor)
				ElevatorP.HandleAsignedOrder(e, btn.Floor, btn.Button, doorStartTimerCh, ps)
			}
		}
	}
}
