package rights

type Service struct {
	//Log *logrus.Logger
}

const (
	AddUser = 1 << (iota)
	EditUser
	DeleteUser
	AddOrganization
	EditOrganization
	DeleteOrganization
	ViewOrganization
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
