package elevator

import "Driver-go/elevio"

type elevatorState int;

const NUMFloors = 4
const numCabButtons = 1
const numHallButtons = 2




const(
	idle elevatorState = 0;
	moving = 1;
	doorOpen = 2;
	error_ = 3;
)

type Elevator struct{
	currentFloor int
	lastFloor int
	cabOrderMatrix  [NUMFloors][numCabButtons]bool
	HallorderMatrix  [NUMFloors][numHallButtons]bool
	dir elevio.MotorDirection
	state elevatorState
	mode elevatorMode
}

type elevatorMode int;
const(
	singleElavator elevatorMode = 0;
	MasterSlave = 1;
)

func init(){
	
}