package ElevatorP

import "Driver-go/elevio"

type elevatorState int;

const NUMFloors = 4
const numCabButtons = 1
const numHallButtons = 2
const timeDoorOpenDuration = 3


//Må ha stor bokstav i starten av ordet i definisjoner av const, type og struct for at de skal bli public variabler 
const(
	Idle elevatorState = 0;
	Moving = 1;
	DoorOpen = 2;
	Error_ = 3;
)

type Elevator struct{
	CurrentFloor int
	LastFloor int
	CabOrderMatrix  [NUMFloors][numCabButtons]bool
	HallorderMatrix  [NUMFloors][numHallButtons]bool
	Dir elevio.MotorDirection
	State elevatorState
	Mode elevatorMode
	obstructed bool
}

type elevatorMode int;
const(
	SingleElavator elevatorMode = 0;
	MasterSlave = 1;
)

func NewElevator() *Elevator{
	e := &Elevator{}
	e.State = Idle
	e.Dir = elevio.MD_Stop
	
	return e

	
}