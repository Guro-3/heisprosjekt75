package ElevatorP

import (
	"heisprosjekt75/types"
	"heisprosjekt75/Driver-go/elevio"
)

func NewElevator(myID string, myIP string) *types.Elevator {

	e := &types.Elevator{
		MyID:   myID,
		ElevIP: myIP,
		State:  types.Idle,
		Dir:    elevio.MD_Stop,
	}

	e.Ps.Role = types.RoleNone
	e.Ps.PrevRole = types.RoleNone

	return e
}