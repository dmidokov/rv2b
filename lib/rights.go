package lib

const (
	AddUser = 1 << (iota)
	_       //EditUser
	DeleteUser
	AddOrganization
	_ //EditOrganization
	_ //DeleteOrganization // todo: добавить проверку этого права при удалении
	ViewOrganization
	ViewBranchList
	AddBranch
	DeleteBranch
	EditUserRights
	EditUserNavigation
	SetUserHotSwitch
	_
)
