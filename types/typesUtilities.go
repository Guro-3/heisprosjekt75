package types

import(
	"heisprosjekt75/Driver-go/elevio"
)

func TypesRoleToString(r ElevatorRole) string {
	switch r {
	case RoleNone:
		return "None"
	case RolePrimary:
		return "Primary"
	case RoleBackup:
		return "Backup"
	case RoleNode:
		return "Node"
	default:
		return "Unknown"
	}
}



func TypesWorldViewToJSON(world map[string]ElevatorStatus) map[string]ElevatorStateJSON {
	result := make(map[string]ElevatorStateJSON)
	for id, e := range world {
		dir := "stop"
		switch e.Direction {
		case elevio.MD_Up:
			dir = "up"
		case elevio.MD_Down:
			dir = "down"
		case elevio.MD_Stop:
			dir = "stop"
		}

		state := "idle"
		switch e.State {
		case Idle:
			state = "idle"
		case Moving:
			state = "moving"
		case DoorOpen:
			state = "doorOpen"
		}

		result[id] = ElevatorStateJSON{
			Behaviour:   state,
			Floor:       e.Floor,
			Direction:   dir,
			CabRequests: e.CabRequests,
		}
	}
	return result
}


func TypesUpdateMyState(e *Elevator) {
	cabs := make([]bool, len(e.CabOrderMatrix))
	copy(cabs, e.CabOrderMatrix[:])

	WorldView[e.MyID] = ElevatorStatus{
		Floor:       e.CurrentFloor,
		Direction:   e.Dir,
		State:       e.State,
		CabRequests: cabs,
	}
}

