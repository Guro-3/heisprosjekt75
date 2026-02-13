package ElevatorP

import (
	"Driver-go/elevio"
	"fmt"
)

func SeCabLight(floor int){
	fmt.Printf("Lights turned on in elevator for %d. floor\n", floor)
	elevio.SetButtonLamp(elevio.BT_Cab, floor, true)
}
func TurnOffCabLight(floor int){
	fmt.Printf("Lights turned of in elevator for %d. floor\n", floor)
	elevio.SetButtonLamp(elevio.BT_Cab, floor, false)
}

func SetHallLight(btn elevio.ButtonType, floor int){
	fmt.Printf("Lights on in hall on floor %d\n", floor)
	elevio.SetButtonLamp(btn, floor, true)
}
func TurnOffHallLight(btn elevio.ButtonType, floor int){
	fmt.Printf("Lights off in hall on floor %d\n", floor)
	elevio.SetButtonLamp(btn, floor, false)
}
func FloorLight(e *Elevator){
	floor := e.CurrentFloor
	elevio.SetFloorIndicator(floor)
}
