package PrimaryHeartbeat

import (
	"heisprosjekt75/types"
	"time"
)

type PrimHeartbeat struct {
	PrimaryID      string
	PrimaryAddrTCP string
}

func SendPrimaryIpId(UDPHeartBeatTx chan PrimHeartbeat, d time.Duration, e *types.Elevator) {
	tic := time.NewTicker(d)
	defer tic.Stop()

	for range tic.C {
		if e.Ps.Role == types.RolePrimary {
			UDPHeartBeatTx <- PrimHeartbeat{e.MyID, e.ElevIP}
		}
	}
}
