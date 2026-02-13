package ElevatorP

import (
	"Driver-go/elevio"
)

func SeCabLight(floor int){
	elevio.SetButtonLamp(elevio.BT_Cab, floor, true)
}
func TurnOffCabLight(floor int){
	elevio.SetButtonLamp(elevio.BT_Cab, floor, false)
}

func SetHallLight(btn elevio.ButtonType, floor int){
	elevio.SetButtonLamp(btn, floor, true)
}
func TurnOffHallLight(btn elevio.ButtonType, floor int){
	elevio.SetButtonLamp(btn, floor, false)
}
func FloorLight(e *Elevator){
	floor := e.CurrentFloor
	elevio.SetFloorIndicator(floor)
}
