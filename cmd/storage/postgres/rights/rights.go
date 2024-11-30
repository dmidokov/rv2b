package rights

import (
	"context"
	"fmt"
	"github.com/dmidokov/rv2/lib/entitie"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"strconv"
)

type Service struct {
	DB  *pgxpool.Pool
	Log *logrus.Logger
}

func New(DB *pgxpool.Pool, Log *logrus.Logger) *Service {
	return &Service{
		DB:  DB,
		Log: Log,
	}
}

func (rg *Service) CheckUserRight(user *entitie.User, right int) bool {
	fmt.Println(user.Rights)
	fmt.Println(right)
	if (user.Rights & right) == right {
		return true
	}
	return false
}

func (rg *Service) GetRightsNamesByIds(ids []int) (*[]entitie.RightNameValue, error) {
	query := `select name, value from rights_names where 1=0 `

	for _, v := range ids {
		query += fmt.Sprintf("OR value=%s", strconv.Itoa(v))
	}

	rows, err := rg.DB.Query(context.Background(), query)
	if err != nil {
		rg.Log.Error()
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
func (rg *Service) GetByUserRights(rightsValue int) (*[]entitie.Right, error) {
	query := "select * from rights_names where value & $1 > 0"

	rows, err := rg.DB.Query(context.Background(), query, rightsValue)
	if err != nil {
		rg.Log.Error()
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

func (rg *Service) GetAvailableEntities(userId int, groupId int) (*[]entitie.Entities, error) {
	query := `select * from rights where user_id = $1 AND entity_group=$2`
	rows, err := rg.DB.Query(context.Background(), query, userId, groupId)

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