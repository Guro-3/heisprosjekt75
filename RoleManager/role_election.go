package RoleManager

import (
	"fmt"
	"heisprosjekt75/Network-go/network/peers"
)

func RoleElection(peers peers.PeerUpdate, MyID string, ps *PeerState) {
	ps.PrevRole = ps.Role

	if len(peers.Peers) == 0 {
		ps.PrimaryID = MyID
	} else {
		ps.PrimaryID = peers.Peers[0]
	}

	if len(peers.Peers) >= 2 {
		ps.BackupID = peers.Peers[1]
	}
	fmt.Printf("DEBUG: MyID=%q PrimaryID=%q BackupID=%q\n",
		MyID, ps.PrimaryID, ps.BackupID)

	switch MyID {
	case ps.PrimaryID:
		ps.Role = RolePrimary
		fmt.Print("my role is Primary\n")
	case ps.BackupID:
		ps.Role = RoleBackup
		fmt.Print("my role is backup\n")
	default:
		ps.Role = RoleNode
		fmt.Print("my role is none\n")
	}
}

