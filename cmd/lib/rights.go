package lib

const (
	AddUser  = 1 << (iota)
	EditUser //EditUser
	DeleteUser
	AddOrganization
	EditOrganization   //EditOrganization
	DeleteOrganization //DeleteOrganization // todo: добавить проверку этого права при удалении
	ViewOrganization
	ViewBranchList
	AddBranch
	DeleteBranch
	EditUserRights
	EditUserNavigation
	EditUserHotSwitch
	HotSwitchToAnotherUser
	ViewUsers
	ViewUserGroups
	EditUserGroups
	DeleteUserGroups
	CreateUserGroup
	AssignUserGroup
	UnassignUserGroup
)
