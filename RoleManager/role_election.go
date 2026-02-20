package RoleManager

import(
	"heisprosjekt75/Network-go/network/peers"
	"fmt"
)

func RoleElection(peers peers.PeerUpdate, MyID string,ps *PeerState){
	ps.PrevRole = ps.Role

	ps.MasterID = ""
    ps.BackupID = ""

	if len(peers.Peers) >=1 {
		ps.MasterID = peers.Peers[0]
	}
	if len(peers.Peers) >= 2 {
		ps. BackupID = peers.Peers[1]
	}
	fmt.Printf("DEBUG: MyID=%q MasterID=%q BackupID=%q\n",
    MyID, ps.MasterID, ps.BackupID)
	
	switch MyID{
	case ps.MasterID:
		ps.Role = RoleMaster
		fmt.Print("my role is Master\n")
	case ps.BackupID:
		ps.Role = RoleBackup
		fmt.Print("my role is backup\n")
	default:
		ps.Role = RoleNode
		fmt.Print("my role is none\n")
	}
}
	
