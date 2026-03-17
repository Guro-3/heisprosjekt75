package schedueler

import (
	"encoding/json"
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/Elevator"
	"heisprosjekt75/types"
	"log"
	"time"
)

func PrimarySchedueler(e *types.Elevator, doorStartTimerCh chan int, excludeFloor int, excludeButton int, excludePeer string) {
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

			if f == excludeFloor && b == excludeButton && owner == excludePeer {
				altOwner := ""
				for id, matrix := range proposedAssignment {
					if id == excludePeer {
						continue
					}
					if matrix[f][b] {
						altOwner = id
						break
					}
				}
				if altOwner != "" {
					owner = altOwner
				}
			}
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

func PrimaryMonitorTick(e *types.Elevator, doorStartTimerCh chan int, d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	var releaseFloor int
	var releaseButton int

	for range ticker.C {
		if e.Ps.Role != types.RolePrimary {
			continue
		}
		needsReschedule := false
		maxOrderAge := 10 * time.Second

		for f := 0; f < types.NumFloors; f++ {
			for b := 0; b < types.NumHallButtons; b++ {
				if !types.FullOrderMatrix[f][b] {
					continue
				}

				t := types.HallOrderTimes[f][b]
				if t.IsZero() {
					types.HallOrderTimes[f][b] = time.Now()
					continue
				}

				if time.Since(t) > maxOrderAge {
					releaseButton = b
					releaseFloor = f

					needsReschedule = true
				}
			}
		}
		if needsReschedule {
			oldOwner := getOwnerOfHallOrder(releaseFloor, releaseButton)
			releaseSpecificHallOrder(releaseFloor, releaseButton)
			PrimarySchedueler(e, doorStartTimerCh, releaseFloor, releaseButton, oldOwner)
		}
	}
}
