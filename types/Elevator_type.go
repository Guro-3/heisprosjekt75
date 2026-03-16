package types

import "heisprosjekt75/Driver-go/elevio"

const NumFloors = 4
const NumCabButtons = 1
const NumHallButtons = 2
const TimeDoorOpenDuration = 3

type ElevatorState int
type ElevatorMode int
type StopCondition int

const (
	Idle ElevatorState = iota
	Moving
	DoorOpen
	Error
)

const (
	SingleElevator ElevatorMode = iota
	PrimaryBackup
)

const (
	UpOrder StopCondition = iota
	DownOrder
	CabOrder
)

type Elevator struct {
	CurrentFloor    int
	LastFloor       int
	CabOrderMatrix  [NumFloors]bool
	HallOrderMatrix [NumFloors][NumHallButtons]bool
	Dir             elevio.MotorDirection
	State           ElevatorState
	Mode            ElevatorMode
	Obstructed      bool
	MyID            string
	StableID        string
	ElevIP          string
	Ps              PeerState
	StopCond        StopCondition
	OrderDir        elevio.MotorDirection
	ClearedRevDir   bool
	Initializing    bool
}
