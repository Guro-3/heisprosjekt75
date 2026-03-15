package RoleManager

import (
	"heisprosjekt75/ElevatorP"
	"heisprosjekt75/Network-go/network/peers"
	"heisprosjekt75/types"
	"log"
)
func contains(peers []string, id string) bool {
	for _, p := range peers {
		if p == id {
			return true
		}
	}
	return false
}

func firstOther(peers []string, avoid string) string {
	for _, p := range peers {
		if p != avoid {
			return p
		}
	}
	return ""
}

func RoleElection(peers peers.PeerUpdate, e *types.Elevator, ps *types.PeerState, doorStartTimerCh chan int) {
	ps.PrevRole = ps.Role

	if len(peers.Peers) == 1 {
		e.Mode = types.SingleElevator
		ps.PrimaryID = ""
		ps.BackupID = ""
		ps.Role = types.RoleNode

		ElevatorP.SingleElevatorOrderRedelegation(e, doorStartTimerCh)
		log.Println("Elevator mode:", e.Mode)
		log.Println("my role is none (single mode)")
		return
	}

	e.Mode = types.PrimaryBackup
	log.Println("Elevator mode:", e.Mode)

	// Behold gammel primary hvis den fortsatt finnes
	if ps.PrimaryID == "" || !contains(peers.Peers, ps.PrimaryID) {
		// Hvis en node nettopp kom tilbake, ikke gjør den til primary med en gang
		if peers.New != "" {
			other := firstOther(peers.Peers, peers.New)
			if other != "" {
				ps.PrimaryID = other
			} else {
				ps.PrimaryID = peers.Peers[0]
			}
		} else {
			ps.PrimaryID = peers.Peers[0]
		}
	}

	ps.BackupID = ""
	for _, id := range peers.Peers {
		if id != ps.PrimaryID {
			ps.BackupID = id
			break
		}
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