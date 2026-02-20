package RoleManager


type elevatorRole int
const(
	RoleMaster elevatorRole = 0;
	RoleBackup = 1;
	RoleNode = 2;
)


type PeerState struct{
	Role elevatorRole
	PrevRole elevatorRole
	MasterID string
	BackupID string
}