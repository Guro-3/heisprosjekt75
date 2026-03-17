package rolemanager

import (
	"heisprosjekt75/Elevator"
	"heisprosjekt75/Network/peers"
	"heisprosjekt75/types"
	"log"
)

func RoleElection(peerUpdate peers.PeerUpdate, e *types.Elevator, doorStartTimerCh chan int) {
	e.Ps.PrevRole = e.Ps.Role
	peerList := peerUpdate.Peers

	if len(peerList) == 1 {
		e.Ps.PrimaryID = e.MyID
		e.Ps.BackupID = ""
		e.Mode = types.SingleElevator
		Elevator.SingleElevatorOrderRedelegation(e, doorStartTimerCh)
	} else {
		e.Mode = types.PrimaryBackup
		e.Ps.PrimaryID = peerList[0]
		e.Ps.BackupID = peerList[1]
	}

	switch e.MyID {
	case e.Ps.PrimaryID:
		e.Ps.Role = types.RolePrimary
		log.Println("my role is Primary")
	case e.Ps.BackupID:
		e.Ps.Role = types.RoleBackup
		log.Println("my role is Backup")
	default:
		e.Ps.Role = types.RoleNode
		log.Println("my role is Node")
	}
}
