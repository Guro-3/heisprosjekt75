package ElevatorP

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/RoleManager"
)


type elevatorState int;



type elevatorMode int;

const NUMFloors = 4
const numCabButtons = 1
const numHallButtons = 2
const timeDoorOpenDuration = 3


const(
	Idle elevatorState = 0;
	Moving = 1;
	DoorOpen = 2;
	Error_ = 3;
)



const(
	SingleElavator elevatorMode = 0;
	MasterSlave = 1;
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
	MyID string
	Ps RoleManager.PeerState
	ElevIP string
}



func NewElevator(myID string, MyIP string) *Elevator{
	e := &Elevator{}
	e.MyID = myID
	e.State = Idle
	e.Dir = elevio.MD_Stop
	e.ElevIP = MyIP
	return e
}

type Heartbeat struct {
	ID string
	Role string
	Port int
	TimeStamp int
}