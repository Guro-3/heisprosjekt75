package ElevatorP

import (
	"heisprosjekt75/Driver-go/elevio"
	sendmessages "heisprosjekt75/Messages/SendMessages"
	"heisprosjekt75/types"
)

func SetCabLight(floor int) {
	elevio.SetButtonLamp(elevio.BT_Cab, floor, true)
}
func TurnOffCabLight(floor int) {
	elevio.SetButtonLamp(elevio.BT_Cab, floor, false)
}

func SetHallLight(btn elevio.ButtonType, floor int) {
	elevio.SetButtonLamp(btn, floor, true)
}
func TurnOffHallLight(btn elevio.ButtonType, floor int) {
	elevio.SetButtonLamp(btn, floor, false)
}
func FloorLight(e *types.Elevator) {
	floor := e.CurrentFloor
	elevio.SetFloorIndicator(floor)
}

func LightInit() {
	elevio.SetDoorOpenLamp(false)
	for f := 0; f < types.NumFloors; f++ {
		for b := 0; b < types.NumHallButtons; b++ {
			TurnOffHallLight(elevio.ButtonType(b), f)
		}

		TurnOffCabLight(f)
	}
}

func SyncHallLight(ps *types.PeerState, e *types.Elevator, world map[string]types.ElevatorStatus){
	for f:= 0; f < types.NumFloors; f++ {
		for b:= 0; b < types.NumHallButtons; b++{
			btn:= elevio.ButtonEvent{
				Floor: f,
				Button: elevio.ButtonType(b),
			}

			if types.FullOrderMatrix[f][b]{
				SetHallLight(btn.Button, btn.Floor)
				sendmessages.SendHallLightOn(ps, e, btn, world)
			}else{
				TurnOffHallLight(btn.Button, btn.Floor)
				sendmessages.SendHallLightOff(ps, e, btn, world)
			}
		}
	}
}