package PrimaryHeartbeat

import (
	//"fmt"
	"heisprosjekt75/ElevatorP"
	"heisprosjekt75/RoleManager"
	"time"
)

type PrimHeartbeat struct {
	PrimaryID      string
	PrimaryAddrTCP string
}

func SendPrimaryIpId(UDPHeartBeatTx chan PrimHeartbeat, d time.Duration, ps *RoleManager.PeerState, e *ElevatorP.Elevator) {
	tic := time.NewTicker(d)
	defer tic.Stop()

	for range tic.C {
		if ps.Role == RoleManager.RolePrimary {
			UDPHeartBeatTx <- PrimHeartbeat{e.MyID, e.ElevIP}
		}
	}
}
