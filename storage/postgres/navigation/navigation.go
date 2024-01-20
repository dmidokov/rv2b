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
	Get(r *http.Request, key string) (interface{}, bool)
}

func New(DB *pgxpool.Pool, CookieStore SessionStorage, Log *logrus.Logger) *Service {
	return &Service{
		DB: DB, CookieStore: CookieStore, Log: Log,
	}
}

func (o *Service) Get(userId int) ([]*e.Navigation, error) {

	query := `
			SELECT 
			    remonttiv2.navigation.navigation_id,
			    remonttiv2.navigation.title,
			    remonttiv2.navigation.tooltip_text,
			    remonttiv2.navigation.navigation_group,
			    remonttiv2.navigation.icon,
			    remonttiv2.navigation.link
			FROM 
			    remonttiv2.navigation, remonttiv2.rights 
			WHERE 
			    remonttiv2.rights.entity_group = $1 AND
			    remonttiv2.rights.user_id = $2 AND
			    remonttiv2.rights.entity_id = remonttiv2.navigation.navigation_id;`

	rows, err := o.DB.Query(context.Background(), query, 1, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRows(rows)

}

func scanRows(rows pgx.Rows) ([]*e.Navigation, error) {
	var result []*e.Navigation
	for rows.Next() {
		item := &e.Navigation{}
		err := rows.Scan(&item.Id, &item.Title, &item.Tooltip, &item.Group, &item.Icon, &item.Link)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}
