package rights

type Service struct{}

const (
	AddUser = 1 << (iota)
	EditUser
	DeleteUser
	AddOrganization
	EditOrganization
	DeleteOrganization // todo: добавить проверку этого права при удалении
	ViewOrganization
	ViewBranchList
	AddBranch
)

func New() *Service {
	return &Service{}
}

func (rg *Service) CheckUserRight(userRights int, right int) bool {
	if (userRights & right) == right {
		return true
	}
	return false
}
