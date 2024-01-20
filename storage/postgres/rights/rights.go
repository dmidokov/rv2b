package rights

import (
	"context"
	"fmt"
	e "github.com/dmidokov/rv2/lib/entitie"
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

func (rg *Service) CheckUserRight(user *e.User, right int) bool {
	if (user.Rights & right) == right {
		return true
	}
	return false
}

func (rg *Service) GetRightsNamesByIds(ids []int) (*[]e.RightNameValue, error) {

	query := `select name, value from remonttiv2.rights_names where 1=0 `

	for _, v := range ids {
		query += fmt.Sprintf("OR value=%s", strconv.Itoa(v))
	}

	rows, err := rg.DB.Query(context.Background(), query)
	if err != nil {
		rg.Log.Error()
		return nil, err
	}

	var result []e.RightNameValue
	for rows.Next() {
		item := e.RightNameValue{}
		err := rows.Scan(&item.Name, &item.Value)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return &result, nil

}

// GetByUserRights возвращает права (индексы, названия и значения) прав пользователя
// пользователям назначены права в таблице users в виде числа, данный метод вернет
// указанные выше значения только для тех прав которые выставлены для пользователя
func (rg *Service) GetByUserRights(rightsValue int) (*[]e.Right, error) {
	query := "select * from remonttiv2.rights_names where value & $1 > 0"

	rows, err := rg.DB.Query(context.Background(), query, rightsValue)
	if err != nil {
		rg.Log.Error()
		return nil, err
	}

	var result []e.Right
	for rows.Next() {
		item := e.Right{}
		err := rows.Scan(&item.Id, &item.Name, &item.Value)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return &result, nil

}
