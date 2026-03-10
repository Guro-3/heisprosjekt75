package RoleManager

import (
	"fmt"
	"heisprosjekt75/Network-go/network/peers"
	"heisprosjekt75/types"
)

func RoleElection(peers peers.PeerUpdate, e *types.Elevator, ps *types.PeerState) {
	ps.PrevRole = ps.Role

	if len(peers.Peers) == 1 {
		ps.PrimaryID = e.MyID
		e.Mode = types.SingleElevator
	} else {
		ps.PrimaryID = peers.Peers[0]
		e.Mode = types.PrimaryBackup
	}

	if len(peers.Peers) >= 2 {
		ps.BackupID = peers.Peers[1]
	}
	fmt.Printf("DEBUG: MyID=%q PrimaryID=%q BackupID=%q\n",
		e.MyID, ps.PrimaryID, ps.BackupID)

	switch e.MyID {
	case ps.PrimaryID:
		ps.Role = types.RolePrimary
		fmt.Print("my role is Primary\n")
	case ps.BackupID:
		ps.Role = types.RoleBackup
		fmt.Print("my role is backup\n")
	default:
		ps.Role = types.RoleNode
		fmt.Print("my role is none\n")
	}
}
