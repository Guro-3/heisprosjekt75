package tcp

import (
	"heisprosjekt75/RoleManager"
	"heisprosjekt75/Driver-go/elevio"
	"net"
)

type msgType struct{
	hallOrder elevio.ButtonType;
	HeartBeat string;
	BackupOrderMatrix [][]int;
	ComplitedOrder string;
}


type nodeHeartbeat struct {
	Role RoleManager.PeerState
	Floor int
	CurrentDir elevio.MotorDirection
}

type MsgData struct {
	MsgType int
	nodeID string
}

var nodeConnMap = make(map[string]net.Conn)