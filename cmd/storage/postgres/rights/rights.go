package rights

import (
	"context"
	"fmt"
	"github.com/dmidokov/rv2/lib/entitie"
	"strconv"
)

func (s *Service) CheckUserRight(user *entitie.User, right int) bool {
	if (user.Rights & right) == right {
		return true
	}
	return false
}

func (s *Service) GetRightsNamesByIds(ids []int) (*[]entitie.RightNameValue, error) {
	query := `select name, value from rights_names where 1=0 `

	for _, v := range ids {
		query += fmt.Sprintf("OR value=%s", strconv.Itoa(v))
	}

	rows, err := s.DB.Query(context.Background(), query)
	if err != nil {
		s.Log.Error()
		return nil, err
	}

	var result []entitie.RightNameValue
	for rows.Next() {
		item := entitie.RightNameValue{}
		err := rows.Scan(&item.Name, &item.Value)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return &result, nil

}

// GetByUserRights возвращает индексы, названия и значения прав пользователя.
// Пользователям назначены права в таблице users в виде числа, данный метод вернет
// указанные выше значения только для тех прав которые выставлены для пользователя
func (s *Service) GetByUserRights(rightsValue int) (*[]entitie.Right, error) {
	query := "select * from rights_names where value & $1 > 0"

	rows, err := s.DB.Query(context.Background(), query, rightsValue)
	if err != nil {
		s.Log.Error()
		return nil, err
	}

	var result []entitie.Right
	for rows.Next() {
		item := entitie.Right{}
		err := rows.Scan(&item.Id, &item.Name, &item.Value)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return &result, nil
}

func (s *Service) GetAvailableEntities(userId int, groupId int) (*[]entitie.Entities, error) {
	query := `select * from rights where user_id = $1 AND entity_group=$2`
	rows, err := s.DB.Query(context.Background(), query, userId, groupId)

	var result []entitie.Entities
	for rows.Next() {
		item := entitie.Entities{}
		err := rows.Scan(&item.UserId, &item.EntityId, &item.GroupId)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *Service) CreateGroup(groupName string, groupRights int64, organizationId int, userId int) error {
	query := "INSERT INTO groups (group_name, group_rights_1, creator_organization_id, creator_id) VALUES ($1, $2, $3, $4)"
	_, err := s.DB.Exec(context.Background(), query, groupName, groupRights, organizationId, userId)
	if err != nil {
		s.Log.Errorf("create group ends with an error: %s", err.Error())
		return err
	}
	return nil
}

func (s *Service) DeleteGroup(groupId int) error {
	query := "DELETE FROM groups user_groups WHERE group_id = $1"
	_, err := s.DB.Exec(context.Background(), query, groupId)
	if err != nil {
		s.Log.Errorf("group delete end's with an error: %s", err.Error())
		return err
	}

	return nil
}

func (s *Service) AssignUserGroup(userId int, groupId int) error {
	query := "INSERT INTO user_groups (user_id, group_id) VALUES ($1, $2)"
	_, err := s.DB.Exec(context.Background(), query, userId, groupId)
	if err != nil {
		s.Log.Errorf("group assign end's with an error: %s", err.Error())
		return err
	}
	return nil
}

func (s *Service) UnassignUserGroup(userId int, groupId int) error {
	query := "DELETE FROM user_groups WHERE user_id=$1 AND group_id=$2"
	_, err := s.DB.Exec(context.Background(), query, userId, groupId)
	if err != nil {
		s.Log.Errorf("group unassign end's with an error: %s", err.Error())
		return err
	}
	return nil
}
