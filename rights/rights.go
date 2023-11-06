package rights

import "github.com/dmidokov/rv2/entitie"

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
	DeleteBranch
)

func New() *Service {
	return &Service{}
}

func (rg *Service) CheckUserRight(user *entitie.User, right int) bool {
	if (user.Rights & right) == right {
		return true
	}
	return false
}
