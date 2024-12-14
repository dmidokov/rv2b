package group

import (
	"github.com/dmidokov/rv2/lib/entitie"
	"net/http"
)

type userProvider interface {
	GetUserIdFromSession(r *http.Request) int
	GetById(userId int) (*entitie.User, error)
}

type rightsProvider interface {
	CheckUserRight(user *entitie.User, right int) bool
	GetGroupsByOrganizationId(organizationId int) ([]entitie.Group, error)
	GetByUserRights(rightsValue int) (*[]entitie.Right, error)
	CreateGroup(groupName string, groupRights int64, organizationId int, userId int) error
	GetGroupByName(groupName string, orgId int) (entitie.Group, error)
	DeleteGroup(groupId int) error
}
