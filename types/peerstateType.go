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
	Role        ElevatorRole
	PrevRole    ElevatorRole
	PrimaryID   string
	PrimaryConn net.Conn
	PrimaryIP   string
	BackupID    string
}

func RoleToString(r ElevatorRole) string {
	switch r {
	case RoleNone:
		return "None"
	case RolePrimary:
		return "Primary"
	case RoleBackup:
		return "Backup"
	case RoleNode:
		return "Node"
	default:
		return "Unknown"
	}
}
