package navigation

import (
	"context"
	e "github.com/dmidokov/rv2/lib/entitie"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Service struct {
	DB          *pgxpool.Pool
	CookieStore SessionStorage
	Log         *logrus.Logger
}

type SessionStorage interface {
	GetByKey(r *http.Request, key string) (interface{}, bool)
}

func New(DB *pgxpool.Pool, CookieStore SessionStorage, Log *logrus.Logger) *Service {
	return &Service{
		DB: DB, CookieStore: CookieStore, Log: Log,
	}
}

func (o *Service) Get(userId int) (*[]e.Navigation, error) {

	query := `
			SELECT 
			    navigation.navigation_id,
			    navigation.title,
			    navigation.tooltip_text,
			    navigation.navigation_group,
			    navigation.icon,
			    navigation.link
			FROM 
			    navigation, rights 
			WHERE 
			    rights.entity_group = $1 AND
			    rights.user_id = $2 AND
			    rights.entity_id = navigation.navigation_id;`

	rows, err := o.DB.Query(context.Background(), query, 1, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRows(rows)

}

func scanRows(rows pgx.Rows) (*[]e.Navigation, error) {
	var result []e.Navigation
	for rows.Next() {
		item := e.Navigation{}
		err := rows.Scan(&item.Id, &item.Title, &item.Tooltip, &item.Group, &item.Icon, &item.Link)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return &result, nil
}

// Set TODO: паренести в работу с entities
func (o *Service) Set(userId int, navigationId int, groupId int) (*e.NavigationAvailable, error) {
	query := "insert into rights (user_id, entity_id, entity_group) values ($1, $2, $3)"
	_, err := o.DB.Exec(context.Background(), query, userId, navigationId, groupId)
	if err != nil {
		return nil, err
	}
	return &e.NavigationAvailable{UserId: userId, EntityId: navigationId, GroupId: 1}, nil
}

// Delete TODO: паренести в работу с entities
func (o *Service) Delete(userId int, navigationId int, groupId int) error {
	query := "delete from rights where user_id = $1 and entity_id=$2 and entity_group=$3"
	_, err := o.DB.Exec(context.Background(), query, userId, navigationId, groupId)
	if err != nil {
		return err
	}
	return nil
}
