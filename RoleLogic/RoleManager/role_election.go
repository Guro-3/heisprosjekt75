package RoleManager

import (
	"heisprosjekt75/ElevatorP"
	"heisprosjekt75/Network-go/network/peers"
	"heisprosjekt75/types"
	"log"
)

func RoleElection(peerUpdate peers.PeerUpdate, e *types.Elevator, ps *types.PeerState, doorStartTimerCh chan int) {
	ps.PrevRole = ps.Role
	peerList := peerUpdate.Peers

	ps.PrimaryID = ""
	ps.PrimaryIP = ""

	if len(peerList) == 1 {
		ps.PrimaryID = e.MyID
		ps.BackupID = ""
		e.Mode = types.SingleElevator
		ElevatorP.SingleElevatorOrderRedelegation(e, doorStartTimerCh)
		log.Println("Elevator mode:", e.Mode)
	} else {
		e.Mode = types.PrimaryBackup
		log.Println("Elevator mode:", e.Mode)

		// Behold eksisterende primary hvis den fortsatt finnes
		ps.PrimaryID = peerList[0]
		ps.BackupID = peerList[1]
	}

	switch e.MyID {
	case ps.PrimaryID:
		ps.Role = types.RolePrimary
		log.Println("my role is Primary")
	case ps.BackupID:
		ps.Role = types.RoleBackup
		log.Println("my role is Backup")
	default:
		ps.Role = types.RoleNode
		log.Println("my role is Node")
	}
}
