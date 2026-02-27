package masterHeartbeat


import(
	"heisprosjekt75/RoleManager"
	"heisprosjekt75/ElevatorP"
	"time"
)

type MstrHeartbeat struct{
	MasterID string;
	MasterAddrTCP string;
}



func SendMasterIpId(UDPHeartBeatTx chan MstrHeartbeat, d time.Duration, ps *RoleManager.PeerState, e *ElevatorP.Elevator){
	tic := time.NewTicker(d)
	defer tic.Stop()

	for range tic.C {
		if ps.Role == RoleManager.RoleMaster{
			UDPHeartBeatTx <- MstrHeartbeat {e.MyID, e.ElevIP}
		}
	}
}