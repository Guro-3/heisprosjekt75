package elevator

import (
	"Driver-go/elevio"
	Driver "Driver-go/elevio"
	"fmt"
)

const NUMFloors = 4
const numButtons = 3


func buttonPressedServiceOrder()

func serviceOrderAtFloor()

type elevatorState int;

const(
	idle elevatorState = 0;
	mooving = 1;
	doorOpen = 2;
	error_ = 3;
)

type Elevator struct{
	currentFloor int
	lastFloor int
	orderMatrix  [NUMFloors][numButtons]bool
	dir elevio.MotorDirection
}

