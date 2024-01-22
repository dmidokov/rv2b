package user

import (
	"context"
	e "github.com/dmidokov/rv2/lib/entitie"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Service struct {
	DB          *pgxpool.Pool
	CookieStore SessionStorage
	Log         *logrus.Logger
}

const (
	DefaultUserRights = 0
	//UserInfoLevel_0   = 0
	InfoFull = 1
)

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

	user, err := scanUser(row)

	if err != nil {
		u.Log.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (u *Service) GetByOrganizationId(userId int) ([]*e.UserShort, error) {

	query := `
			select distinct 
    			user_id, user_name, create_time, update_time, user_type 
			from 
			    remonttiv2.users as users,
				remonttiv2.users_create_relations as relations
			where 
			    (
			        (users.user_id = relations.created_id AND relations.creator_id = $1) OR 
			        users.user_id = $1
			    )
			order by user_id;
`

	rows, err := u.DB.Query(context.Background(), query, userId)
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
		err := rows.Scan(&item.Id, &item.UserName, &item.CreateTime, &item.UpdateTime, &item.Type)
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

	user, err := scanUser(row)

	if err != nil {
		u.Log.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (u *Service) Create(user *e.User) (int, error) {
	// Обновляя поля в этом запросе, обновить и update
	query := `
		INSERT INTO remonttiv2.users
			(user_name, user_password, create_time, update_time, organization_id, rights_1, actions_code, user_type, start_page) 
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING user_id;`

	row := u.DB.QueryRow(
		context.Background(),
		query,
		user.UserName,
		user.Password,
		user.CreateTime,
		user.UpdateTime,
		user.OrganizationId,
		user.Rights,
		user.ActionCode,
		user.Type,
		user.StartPage,
	)

	var id int
	err := row.Scan(&id)

	if err != nil {
		return id, err
	}

	return id, nil
}

func (u *Service) GetIcon(userId int) *e.UserIcon {
	user, err := u.GetById(userId)
	if err != nil {
		return &e.UserIcon{ImageName: ""}
	}

	return &e.UserIcon{ImageName: user.Icon}
}

func scanUser(row pgx.Row) (*e.User, error) {
	var user = &e.User{}

	err := row.Scan(
		&user.Id,
		&user.OrganizationId,
		&user.UserName,
		&user.Password,
		&user.ActionCode,
		&user.Rights,
		&user.CreateTime,
		&user.UpdateTime,
		&user.Icon,
		&user.Type,
		&user.StartPage,
	)

	if err != nil {
		return nil, err
	}

	return user, nil

}

func (u *Service) SetIcon(userId int, iconLink string) error {

	query := `
		UPDATE remonttiv2.users SET account_icon = $1 WHERE user_id = $2`

	_, err := u.DB.Exec(context.Background(), query, iconLink, userId)
	if err != nil {
		logrus.Warning(err.Error())
	}

	return nil
}

func (u *Service) SetUserCreateRelations(creatorId int, createdId int) error {
	query := `
		INSERT INTO remonttiv2.users_create_relations (creator_id, created_id) VALUES($1, $2)`
	_, err := u.DB.Exec(context.Background(), query, creatorId, createdId)
	if err != nil {
		logrus.Warning(err.Error())
	}

	return nil
}

func (u *Service) GetChild(userId int) ([]*e.UserShort, error) {

	query := `
			select distinct 
    			user_id, user_name, create_time, update_time, user_type 
			from 
			    remonttiv2.users as users,
				remonttiv2.users_create_relations as relations
			where 
			    (
			        (users.user_id = relations.created_id AND relations.creator_id = $1)
			    );
`

	rows, err := u.DB.Query(context.Background(), query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanRows(rows)

}

func (u *Service) GetInfo(userId int, infoLevel int) (*e.UserInfoFull, error) {
	userFull := e.UserInfoFull{}
	switch infoLevel {
	case InfoFull:

		user, err := u.GetById(userId)
		if err != nil {
			return nil, err
		}

		var userRightsList []int

		for i := 1; i <= user.Rights; i = i << 1 {
			if (user.Rights & i) > 0 {
				userRightsList = append(userRightsList, i)
			}
		}

		//TODO: вернуть еще параметров
		userFull.UserName = user.UserName
		userFull.Id = user.Id
		userFull.OrganizationId = user.OrganizationId
		userFull.Password = user.Password
		userFull.ActionCode = user.ActionCode
		userFull.UserRights = userRightsList
		userFull.CreateTime = user.CreateTime
		userFull.UpdateTime = user.UpdateTime
		userFull.Icon = user.Icon
		userFull.Type = user.Type
		userFull.StartPage = user.StartPage
	}
	return &userFull, nil
}

func (u *Service) UpdateUser(user *e.User) (*e.User, error) {
	// Обновляя поля в этом запросе, обновить и create
	query := `
		update remonttiv2.users set 
		    user_name = $1, 
		    user_password = $2, 
		    update_time = $3, 
		    organization_id = $4, 
		    rights_1 = $5, 
		    actions_code = $6, 
		    user_type = $7, 
		    start_page = $8
		where user_id = $9
	`
	p, err := u.DB.Exec(
		context.Background(),
		query,
		user.UserName,
		user.Password,
		time.Now().Unix(),
		user.OrganizationId,
		user.Rights,
		user.ActionCode,
		user.Type,
		user.StartPage,
		user.Id,
	)

	logrus.Info(p.String())

	if err != nil {
		return nil, err
	}

	return u.GetById(user.Id)
}

func (u *Service) GetParentId(userId int) (int, error) {
	query := `select creator_id from remonttiv2.users_create_relations where created_id=$1`
	result := u.DB.QueryRow(context.Background(), query, userId)

	var id *int

	err := result.Scan(&id)
	if err != nil {
		return 0, err
	}

	return *id, nil
}

func (u *Service) GetParentUser(userId int) (*e.User, error) {
	query := `
			select distinct * 
			from 
				remonttiv2.users as users, remonttiv2.users_create_relations as relations
			where 
				users.user_id=relations.creator_id AND
				relations.created_id = $1
    `

	row := u.DB.QueryRow(context.Background(), query, userId)

	user, err := scanUser(row)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *Service) SetHotSwitchRelation(fromId, toId int) error {
	query := `insert into remonttiv2.hot_switch_relations (from_user, to_user) values ($1, $2)`
	_, err := u.DB.Exec(context.Background(), query, fromId, toId)
	if err != nil {
		return err
	}
	return nil
}
