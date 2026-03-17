package types

import "net"

type ElevatorRole int

const (
	RoleNone ElevatorRole = iota
	RolePrimary
	RoleBackup
	RoleNode
)

type PeerState struct {
	Role            ElevatorRole
	PrevRole        ElevatorRole
	PrimaryID       string
	PrimaryConn     net.Conn
	PrimaryListener net.Listener
	PrimaryIP       string
	BackupID        string
}
