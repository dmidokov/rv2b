package rights

import (
	"context"
	"github.com/dmidokov/rv2/lib/entitie"
)

func (s *Service) GetGroupsByOrganizationId(organizationId int) ([]entitie.Group, error) {
	var out []entitie.Group

	query := `select groups.group_id, groups.creator_organization_id, groups.group_name, groups.group_rights_1, groups.creator_id, users.user_name from groups, users where creator_organization_id = $1 and creator_id = user_id`
	rows, err := s.DB.Query(context.Background(), query, organizationId)
	if err != nil {
		s.Log.Errorf("DB error: %s", err.Error())
		return out, err
	}

	for rows.Next() {
		g := entitie.Group{}
		if err := rows.Scan(&g.GroupId, &g.CreatorOrganizationId, &g.GroupName, &g.GroupRights1, &g.CreatorId, &g.CreatorName); err != nil {
			s.Log.Errorf("group parse error: %s", err.Error())
		} else {
			out = append(out, g)
		}
	}

	return out, nil

}

func (s *Service) GetGroupByName(groupName string, orgId int) (entitie.Group, error) {

	group := entitie.Group{}

	query := `select * from groups where group_name=$1 and creator_organization_id=$2 LIMIT 1`
	row := s.DB.QueryRow(context.Background(), query, groupName, orgId)

	err := row.Scan(&group.GroupId, &group.CreatorOrganizationId, &group.GroupName, &group.GroupRights1, &group.CreatorId)
	if err != nil {
		s.Log.Errorf("DB error: %s", err.Error())
		return group, err
	}

	return group, nil

}

func (s *Service) GetUserGroupsWithName(userId int) ([]entitie.GroupNameAndIds, error) {
	var group []entitie.GroupNameAndIds

	query := "SELECT user_groups.user_id, user_groups.group_id, groups.group_name FROM user_groups, groups WHERE user_groups.user_id=$1 AND groups.group_id = user_groups.group_id"
	rows, err := s.DB.Query(context.Background(), query, userId)
	if err != nil {
		s.Log.Errorf("DB error: %s", err.Error())
		return group, err
	}

	for rows.Next() {
		g := entitie.GroupNameAndIds{}

		if err := rows.Scan(&g.UserId, &g.GroupId, &g.GroupName); err != nil {
			s.Log.Errorf("group parse error: %s", err.Error())
		} else {
			group = append(group, g)
		}
	}
	return group, nil
}
