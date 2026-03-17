package sendworldview

import (
	"heisprosjekt75/types"
	"time"
)

type PrimaryIPID struct {
	PrimaryID      string
	PrimaryAddrTCP string
}

func SendPrimaryIpId(UDPPrimaryIPIDTx chan PrimaryIPID, d time.Duration, e *types.Elevator) {
	tic := time.NewTicker(d)
	defer tic.Stop()

	for range tic.C {
		if e.Ps.Role == types.RolePrimary {
			UDPPrimaryIPIDTx <- PrimaryIPID{e.MyID, e.ElevIP}
		}
	}
}
