package RoleManager

import (
	"heisprosjekt75/ElevatorP"
	"heisprosjekt75/Network-go/network/peers"
	"heisprosjekt75/types"
	"log"
)

func RoleElection(peers peers.PeerUpdate, e *types.Elevator, ps *types.PeerState, doorStartTimerCh chan int) {
	ps.PrevRole = ps.Role

	if len(peers.Peers) == 1 {
		//	ps.PrimaryID = e.MyID
		e.Mode = types.SingleElevator
		ElevatorP.SingleElevatorOrderRedelegation(e, doorStartTimerCh)
		log.Println("Elevator mode: ", e.Mode)
	} else {
		ps.PrimaryID = peers.Peers[0]
		e.Mode = types.PrimaryBackup
		log.Println("Elevator mode: ", e.Mode)
	}

	if len(peers.Peers) >= 2 {
		ps.BackupID = peers.Peers[1]
	}

	switch e.MyID {
	case ps.PrimaryID:
		ps.Role = types.RolePrimary
		log.Println("my role is Primary")
	case ps.BackupID:
		ps.Role = types.RoleBackup
		log.Println("my role is backup")
	default:
		ps.Role = types.RoleNode
		log.Println("my role is none")
	}
}
