package ElevatorP

import (
	"heisprosjekt75/Driver-go/elevio"
	"heisprosjekt75/types"
	"log"
)

func NewElevator(myID string, myIP string) *types.Elevator {
	LightInit()

	e := &types.Elevator{
		MyID:   myID,
		ElevIP: myIP,
		State:  types.Idle,
		Dir:    elevio.MD_Stop,
		Mode:   types.SingleElevator,
	}
	log.Println("Elevator mode: ", e.Mode)
	e.Ps.Role = types.RoleNone
	e.Ps.PrevRole = types.RoleNone

	
	return e
}
