package Networkloss

import (
	"heisprosjekt75/Network-go/network/localip"
	"time"
	"heisprosjekt75/types"
	"fmt"
	"heisprosjekt75/ElevatorP"
)
func checkConnection() bool {
	_, err := localip.LocalIP()
	return err == nil
}


func PollNetworkConnection(networkConnection chan<- bool) {
	last := checkConnection()
	networkConnection <- last

	for {
		current := checkConnection()
		if current != last {
			networkConnection <- current
			last = current
		}
		time.Sleep(2 * time.Second)
	}
}



// ENDRET
func HandleNetworkLost(ps *types.PeerState, e *types.Elevator, doorStartTimerCh chan int) {
	fmt.Println("NETWORK LOST -> går til single elevator mode")

	if ps.PrimaryConn != nil {
		_ = ps.PrimaryConn.Close()
		ps.PrimaryConn = nil
	}

	ps.PrimaryID = ""
	ps.BackupID = ""

	ps.PrevRole = ps.Role
	ps.Role = types.RolePrimary

	e.Mode = types.SingleElevator
	ElevatorP.SingleElevatorOrderRedelegation(e, doorStartTimerCh)
}

func HandleNetworkRestored(ps *types.PeerState, e *types.Elevator) {
	fmt.Println("NETWORK RESTORED -> venter på stabil peer discovery")

	if ps.PrimaryConn != nil {
		_ = ps.PrimaryConn.Close()
		ps.PrimaryConn = nil
	}

	ps.PrimaryID = ""
	ps.BackupID = ""

	ps.Role = types.RoleNode
	e.Mode = types.PrimaryBackup
}

