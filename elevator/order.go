package elevator

import "Driver-go/elevio"


func addOrder()

func cabOrdersHere(elevator *Elevator){
	return elevator.orderMatrix[elevator.currentFloor][elevio.BT_Cab]
}




// sjekke om det er en hallOrder opp i etasjen vi er i elevator er selve heisen, innholder, etasje vi er i og direction vi beveger oss i og order matrix
func hallOrderUPHere(elevator *Elevator)bool{
	return elevator.orderMatrix[elevator.currentFloor][elevio.BT_HallUp]
	
}

// sjekke om det er en hallOrder ned i etasjen vi er i elevator er selve heisen, innholder, etasje vi er i og direction vi beveger oss i og order matrix
func hallOrderDownHere(elevator *Elevator)bool{
	return elevator.orderMatrix[elevator.currentFloor][elevio.BT_HallDown]
}
	
// sjekke om det er en ordre over der vi er i etasjen vi er i elevator er selve heisen, innholder, etasje vi er i og direction vi beveger oss i og order matrix
func orderBelow(elevator *Elevator)bool{
	for f := elevator.currentFloor - 1; f>= 0; f--{
		for b := 0; b < numButtons; b++{
			if elevator.orderMatrix[f][b]{
				return true
			}
		}
	}
	return false
}
// sjekke om det er en ordre under der vi er der vi er i etasjen vi er i elevator er selve heisen, innholder, etasje vi er i og direction vi beveger oss i og order matrix
func orderAbove(elevator *Elevator)bool{
	for f := elevator.currentFloor + 1; f < NUMFloors; f++{
		for b := 0; b < numButtons; b++{
			if elevator.orderMatrix[f][b]{
				return true
			}
		}
	}
	return false
}

func chooseDirection()
func shouldStop()
func clearFoor()

