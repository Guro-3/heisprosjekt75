package schedueler

import (
	"encoding/json"
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/ElevatorP"
	"heisprosjekt75/Network-go/network/tcp"
	"heisprosjekt75/types"
	"log"
	"os/exec"
	"time"
)

const (
	OrderTimeout = 5 * time.Second
	TickerPeriod = 500 * time.Millisecond
)

type HallOrderInfo struct {
	AddedTime  time.Time
	Floor      int
	Button     elevio.ButtonType
	AssignedTo string
	Completed  bool
}

var ActiveOrders [types.NumFloors][types.NumHallButtons]HallOrderInfo

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

func chooseOwner(floor int, button int, proposedAssignment map[string]types.HAllAssignment, finalAssignment map[string]types.HAllAssignment, exclude string) string {
	owner := ""
	other := 1 - button

	for id, matrix := range types.CurrentAssignment {
		if id == exclude {
			continue
		}
		_, alive := types.WorldView[id]
		if alive && matrix[floor][button] && !finalAssignment[id][floor][other] {
			owner = id
			break
		}
	}

	if owner == "" {
		for id, matrix := range proposedAssignment {
			if id == exclude {
				continue
			}
			_, alive := types.WorldView[id]
			if alive && matrix[floor][button] && !finalAssignment[id][floor][other] {
				owner = id
				break
			}
		}
	}

	if owner == "" {
		for id := range types.WorldView {
			if id == exclude {
				continue
			}
			_, alive := types.WorldView[id]
			if alive && !finalAssignment[id][floor][other] {
				owner = id
				break
			}
		}
	}

	return owner
}

func DelegateOrders(receiverID string, ps *types.PeerState, e *types.Elevator, btn elevio.ButtonEvent) {
	messageData := tcp.HallOrderMessage{Floor: btn.Floor, Button: btn.Button}
	buttonMessage := tcp.Message{
		Type:        tcp.MsgHallOrder,
		NodeID:      e.MyID,
		MessageData: messageData,
	}
	tcp.SendTCP(receiverID, buttonMessage, ps)
}

func checkAndRedelegate(e *types.Elevator, ps *types.PeerState, doorStartTimerCh chan int) {

	for f := 0; f < types.NumFloors; f++ {
		for b := 0; b < types.NumHallButtons; b++ {
			info := ActiveOrders[f][b]
			if !info.Completed && !info.AddedTime.IsZero() && time.Since(info.AddedTime) > OrderTimeout {
				log.Printf("Redelegating order floor:%d button:%d from %s\n", f, b, info.AssignedTo)
				reassignOrder(f, elevio.ButtonType(b), info.AssignedTo, e, ps, doorStartTimerCh)
			}
		}
	}
}

func reassignOrder(floor int, button elevio.ButtonType, oldOwner string, e *types.Elevator, ps *types.PeerState, doorStartTimerCh chan int) {
	finalAssignment := make(map[string]types.HAllAssignment)
	for id := range types.WorldView {
		finalAssignment[id] = types.HAllAssignment{}
	}

	newOwner := chooseOwner(floor, int(button), types.CurrentAssignment, finalAssignment, oldOwner)
	if newOwner != "" {
		ActiveOrders[floor][button].AssignedTo = newOwner
		btn := elevio.ButtonEvent{Floor: floor, Button: button}

		if newOwner == e.MyID {
			ElevatorP.HandleAsignedOrder(e, floor, button, doorStartTimerCh, ps)
		} else {
			DelegateOrders(newOwner, ps, e, btn)
		}
	}
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

			owner := chooseOwner(f, b, proposedAssignment, finalAssignment, "")

			if owner != "" {
				m := finalAssignment[owner]
				m[f][b] = true
				finalAssignment[owner] = m

				btn := elevio.ButtonEvent{
					Floor:  f,
					Button: elevio.ButtonType(b),
				}

				// Oppdater aktive ordrer
				ActiveOrders[f][b] = HallOrderInfo{
					AddedTime:  time.Now(),
					Floor:      f,
					Button:     elevio.ButtonType(b),
					AssignedTo: owner,
					Completed:  false,
				}

				if owner == e.MyID {
					ElevatorP.HandleAsignedOrder(e, f, elevio.ButtonType(b), doorStartTimerCh, ps)
				} else {
					DelegateOrders(owner, ps, e, btn)
				}
			}
		}
	}

	types.CurrentAssignment = finalAssignment
}

func StartRedelegationLoop(e *types.Elevator, ps *types.PeerState, doorStartTimerCh chan int) {
	go func() {
		ticker := time.NewTicker(TickerPeriod)
		for range ticker.C {
			checkAndRedelegate(e, ps, doorStartTimerCh)
		}
	}()
}
