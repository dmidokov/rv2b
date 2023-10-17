package user

import (
	"context"
	e "github.com/dmidokov/rv2/entitie"
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

const DefaultUserRights = 0

type SessionStorage interface {
	Save(r *http.Request, w http.ResponseWriter, data map[string]interface{}) bool
	Get(r *http.Request, key string) (interface{}, bool)
}

func New(DB *pgxpool.Pool, CookieStore SessionStorage, Log *logrus.Logger) *Service {
	return &Service{
		DB:          DB,
		CookieStore: CookieStore,
		Log:         Log,
	}
}

func (u *Service) GetUserByLoginAndOrganization(login string, organizationId int) (*e.User, error) {

	var user = &e.User{}

	query := `select * from remonttiv2.users where user_name = $1 and organization_id=$2;`

	row := u.DB.QueryRow(context.Background(), query, login, organizationId)

	err := row.Scan(
		&user.Id,
		&user.OrganizationId,
		&user.UserName,
		&user.Password,
		&user.ActionCode,
		&user.Rights,
		&user.CreateTime,
		&user.UpdateTime,
	)

	if err != nil {
		u.Log.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (u *Service) GetByOrganizationId(orgId int) ([]*e.UserShort, error) {

	query := `
			select 
    			user_id, user_name, create_time, update_time 
			from 
			    remonttiv2.users 
			where 
			    organization_id = $1;
`

	rows, err := u.DB.Query(context.Background(), query, orgId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanRows(rows)

}

// GetUserIdFromSession возвращает айди клиента в сессии или ноль
func (u *Service) GetUserIdFromSession(r *http.Request) int {
	log := u.Log

	if auth, ok := u.CookieStore.Get(r, "authenticated"); !ok || !auth.(bool) {
		log.Warning("User is not authorized")
		return 0
	}

	if userId, ok := u.CookieStore.Get(r, "userid"); !ok {
		log.Warning("User is not authorized")
		return 0
	} else {
		return userId.(int)
	}
}

func (u *Service) IsAuthorized(r *http.Request) bool {
	if u.GetUserIdFromSession(r) != 0 {
		return true
	}
	return false
}

// GetOrganizationIdFromSession возвращает айди организации клиента в сессии или ноль
func (u *Service) GetOrganizationIdFromSession(r *http.Request) int {
	log := u.Log

	if auth, ok := u.CookieStore.Get(r, "authenticated"); !ok || !auth.(bool) {
		log.Warning("User is not authorized")
		return 0
	}

	if orgId, ok := u.CookieStore.Get(r, "organizationid"); !ok {
		log.Warning("User is not authorized")
		return 0
	} else {
		return orgId.(int)
	}
}

func scanRows(rows pgx.Rows) ([]*e.UserShort, error) {
	var result []*e.UserShort
	for rows.Next() {
		item := &e.UserShort{}
		err := rows.Scan(&item.Id, &item.UserName, &item.CreateTime, &item.UpdateTime)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

func (u *Service) Delete(userId int) error {
	query := `
		DELETE FROM remonttiv2.users 
		WHERE user_id=$1`

	tag, err := u.DB.Exec(context.Background(), query, userId)
	if err != nil {
		return err
	}
	u.Log.Info("User deleted: ", tag.RowsAffected())

	return nil
}

func (u *Service) GetById(userId int) (*e.User, error) {
	var user = &e.User{}

	query := `select * from  remonttiv2.users where user_id = $1;`

	row := u.DB.QueryRow(context.Background(), query, userId)

	err := row.Scan(
		&user.Id,
		&user.OrganizationId,
		&user.UserName,
		&user.Password,
		&user.ActionCode,
		&user.Rights,
		&user.CreateTime,
		&user.UpdateTime,
	)

	if err != nil {
		u.Log.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (u *Service) Create(user *e.User) error {
	query := `
		INSERT INTO remonttiv2.users
			(user_name, user_password, create_time, update_time, organization_id, rights_1, actions_code) 
		VALUES
			($1, $2, $3, $4, $5, $6, $7);`

	tag, err := u.DB.Exec(context.Background(), query, user.UserName, user.Password, user.CreateTime, user.UpdateTime, user.OrganizationId, user.Rights, user.ActionCode)
	if err != nil {
		return err
	}
	u.Log.Info("Создано пользователей:", tag.RowsAffected())

	return nil
}
