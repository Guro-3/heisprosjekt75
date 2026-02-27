package RoleManager

import(
	"net"
)


type elevatorRole int
const(
	RolePrimary elevatorRole = 0;
	RoleBackup = 1;
	RoleNode = 2;
)


type PeerState struct{
	Role elevatorRole
	PrevRole elevatorRole
	PrimaryID string
	PrimaryConn net.Conn
	PrimaryIP string
	BackupID string
	
}