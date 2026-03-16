package Elevator

import(
	"heisprosjekt75/types"
	"heisprosjekt75/Driver-go/elevio"
)

func checkCabOrdersHere(e *types.Elevator) bool {
	return e.CabOrderMatrix[e.CurrentFloor]
}

func checkHallOrderUpHere(e *types.Elevator) bool {
	return e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallUp]
}

func checkHallOrderDownHere(e *types.Elevator) bool {
	return e.HallOrderMatrix[e.CurrentFloor][elevio.BT_HallDown]
}

func checkCabOrdersBelow(e *types.Elevator) bool {
	for f := e.CurrentFloor - 1; f >= 0; f-- {
		if e.CabOrderMatrix[f] {
			return true
		}
	}
	return false
}
func checkCabOrdersAbove(e *types.Elevator) bool {
	for f := e.CurrentFloor + 1; f < types.NumFloors; f++ {
		if e.CabOrderMatrix[f] {
			return true
		}
	}
	return false
}

func checkOrderBelow(e *types.Elevator) bool {
	for f := e.CurrentFloor - 1; f >= 0; f-- {
		for b := 0; b < types.NumHallButtons; b++ {
			if e.HallOrderMatrix[f][b] {
				return true
			}
		}

		if e.CabOrderMatrix[f] {
			return true
		}
	}
	return false
}

func checkOrderAbove(e *types.Elevator) bool {
	for f := e.CurrentFloor + 1; f < types.NumFloors; f++ {
		for b := 0; b < types.NumHallButtons; b++ {
			if e.HallOrderMatrix[f][b] {
				return true
			}
		}
		if e.CabOrderMatrix[f] {
			return true
		}
	}
	return false
}