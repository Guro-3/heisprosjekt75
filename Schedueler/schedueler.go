package schedueler

import (
	"encoding/json"
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/Elevator"
	"heisprosjekt75/types"
	"log"
)




func PrimarySchedueler(e *types.Elevator, doorStartTimerCh chan int) {
	types.TypesUpdateMyState(e)

	hallRequests := make([][types.NumHallButtons]bool, types.NumFloors)

	for f := 0; f < types.NumFloors; f++ {
		hallRequests[f] = [types.NumHallButtons]bool{
			types.FullOrderMatrix[f][elevio.BT_HallUp],
			types.FullOrderMatrix[f][elevio.BT_HallDown],
		}
	}

	input := types.HRAInput{
		HallRequests: hallRequests,
		States:       types.TypesWorldViewToJSON(types.WorldView),
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		log.Println("JSON Marshal error in PrimarySchedueler:", err)
		return
	}

	proposal, err := assignHallRequests(jsonBytes)
	if err != nil {
		log.Println("AssignHallRequests error:", err)
		return
	}

	proposedAssignment := make(map[string]types.HallAssignment)

	for id, matrix := range proposal {
		proposedAssignment[id] = hallAssignmentConvert(matrix)
	}

	finalAssignment := make(map[string]types.HallAssignment)

	for id := range types.WorldView {
		finalAssignment[id] = types.HallAssignment{}
	}

	for f := 0; f < types.NumFloors; f++ {
		for b := 0; b < types.NumHallButtons; b++ {

			if !types.FullOrderMatrix[f][b] {
				continue
			}
			owner := chooseOwner(f, b, proposedAssignment, finalAssignment)
			if owner != "" {
				order := finalAssignment[owner]
				order[f][b] = true
				finalAssignment[owner] = order
			}
		}
	}
	for id, newOrderMatrix := range finalAssignment {
		oldOrderMatrix := types.CurrentAssignment[id]
		for f := 0; f < types.NumFloors; f++ {
			for b := 0; b < types.NumHallButtons; b++ {

				if !oldOrderMatrix[f][b] && newOrderMatrix[f][b] {
					btn := elevio.ButtonEvent{
						Floor:  f,
						Button: elevio.ButtonType(b),
					}
					if id == e.MyID {
						Elevator.HandleAssignedOrder(e, f, elevio.ButtonType(b), doorStartTimerCh)
					} else {
						delegateOrders(id, e, btn)
					}
				}
			}
		}
	}
	types.CurrentAssignment = finalAssignment
}
