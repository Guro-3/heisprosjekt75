package Elevator

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/types"
)

func InitNewElevator(myID string, myIP string) *types.Elevator {
	LightInit()

	e := &types.Elevator{
		MyID:   myID,
		ElevIP: myIP,
		State:  types.Idle,
		Dir:    elevio.MD_Stop,
		Mode:   types.SingleElevator,
	}
	e.Ps.Role = types.RoleNone
	e.Ps.PrevRole = types.RoleNone

	return e
}

func InitBetweenFloor(e *types.Elevator) {
	e.Initializing = true
	elevio.SetMotorDirection(elevio.MD_Down)
	e.Dir = elevio.MD_Down
	e.State = types.Moving
}
