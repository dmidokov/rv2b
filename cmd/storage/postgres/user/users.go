package user

import (
	"context"
	"fmt"
	"github.com/dmidokov/rv2/lib/entitie"
	"github.com/dmidokov/rv2/session/cookie"
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
	GetByKey(r *http.Request, key string) (interface{}, bool)
}

func New(DB *pgxpool.Pool, CookieStore SessionStorage, Log *logrus.Logger) *Service {
	return &Service{
		DB:          DB,
		CookieStore: CookieStore,
		Log:         Log,
	}
}

func (u *Service) GetUserByLoginAndOrganization(login string, organizationId int) (*entitie.User, error) {

	var user = &entitie.User{}

	query := `select * from users where user_name = $1 and organization_id=$2;`

	row := u.DB.QueryRow(context.Background(), query, login, organizationId)

	user, err := scanUser(row)

	if err != nil {
		u.Log.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (u *Service) GetByOrganizationId(userId int) ([]*entitie.UserShort, error) {

	query := `
			select distinct 
    			user_id, user_name, create_time, update_time, user_type 
			from 
			    users as users,
				users_create_relations as relations
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

	if auth, ok := u.CookieStore.GetByKey(r, cookie.Authenticated); !ok || !auth.(bool) {
		log.Warning("User is not authorized")
		return 0
	}

	if userId, ok := u.CookieStore.GetByKey(r, cookie.UserId); !ok {
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

	if auth, ok := u.CookieStore.GetByKey(r, cookie.Authenticated); !ok || !auth.(bool) {
		log.Warning("User is not authorized")
		return 0
	}

	if orgId, ok := u.CookieStore.GetByKey(r, cookie.OrganizationId); !ok {
		log.Warning("User is not authorized")
		return 0
	} else {
		return orgId.(int)
	}
}

func scanRows(rows pgx.Rows) ([]*entitie.UserShort, error) {
	var result []*entitie.UserShort
	for rows.Next() {
		item := &entitie.UserShort{}
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
		DELETE FROM users 
		WHERE user_id=$1`

	tag, err := u.DB.Exec(context.Background(), query, userId)
	if err != nil {
		return err
	}
	u.Log.Info("User deleted: ", tag.RowsAffected())

	return nil
}

func (u *Service) GetById(userId int) (*entitie.User, error) {
	var user = &entitie.User{}

	query := `select * from  users where user_id = $1;`

	row := u.DB.QueryRow(context.Background(), query, userId)

	user, err := scanUser(row)

	if err != nil {
		u.Log.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (u *Service) Create(user *entitie.User) (int, error) {
	// Обновляя поля в этом запросе, обновить и update
	query := `
		INSERT INTO users
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

func (u *Service) GetIcon(userId int) *entitie.UserIcon {
	user, err := u.GetById(userId)
	if err != nil {
		return &entitie.UserIcon{ImageName: ""}
	}

	return &entitie.UserIcon{ImageName: user.Icon}
}

func scanUser(row pgx.Row) (*entitie.User, error) {
	var user = &entitie.User{}

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
		UPDATE users SET account_icon = $1 WHERE user_id = $2`

	_, err := u.DB.Exec(context.Background(), query, iconLink, userId)
	if err != nil {
		logrus.Warning(err.Error())
	}

	return nil
}

func (u *Service) SetUserCreateRelations(creatorId int, createdId int) error {
	query := `
		INSERT INTO users_create_relations (creator_id, created_id) VALUES($1, $2)`
	_, err := u.DB.Exec(context.Background(), query, creatorId, createdId)
	if err != nil {
		logrus.Warning(err.Error())
	}

	return nil
}

func (u *Service) GetChild(userId int) ([]*entitie.UserShort, error) {

	query := `
			select distinct 
    			user_id, user_name, create_time, update_time, user_type 
			from 
			    users as users,
				users_create_relations as relations
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

func (u *Service) GetInfo(userId int, infoLevel int) (*entitie.UserInfoFull, error) {
	userFull := entitie.UserInfoFull{}
	switch infoLevel {
	case InfoFull:

		user, err := u.GetById(userId)
		if err != nil {
			return nil, err
		}

		var userRightsList []int64
		var i int64

		for i = 1; i <= user.Rights; i = i << 1 {
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

func (u *Service) UpdateUser(user *entitie.User) (*entitie.User, error) {
	// Обновляя поля в этом запросе, обновить и create
	query := `
		update users set 
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
	query := `select creator_id from users_create_relations where created_id=$1`
	result := u.DB.QueryRow(context.Background(), query, userId)

	var id *int

	err := result.Scan(&id)
	if err != nil {
		return 0, err
	}

	return *id, nil
}

func (u *Service) GetParentUser(userId int) (*entitie.User, error) {
	query := `
			select distinct * 
			from 
				users as users, users_create_relations as relations
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
	query := `insert into hot_switch_relations (from_user, to_user) values ($1, $2)`
	_, err := u.DB.Exec(context.Background(), query, fromId, toId)
	if err != nil {
		return err
	}
	return nil
}

func (u *Service) RemoveHotSwitchRelation(fromId, toId int) error {
	query := `delete from hot_switch_relations where from_user=$1 and  to_user=$2`
	_, err := u.DB.Exec(context.Background(), query, fromId, toId)
	if err != nil {
		return err
	}
	return nil
}

// TODO: две функции ниже делают одно и тоже, надо объединить

func (u *Service) GetHotSwitch(userId int) ([]*entitie.UserShort, error) {
	query := `select users.user_id, users.user_name, users.create_time, users.update_time, users.user_type from users as users, hot_switch_relations as switch where users.user_id = switch.to_user and switch.from_user=$1`
	rows, err := u.DB.Query(context.Background(), query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanRows(rows)
}

func (u *Service) GetUsersToSwitch(userId int) ([]*entitie.UserSwitcher, error) {
	query := `select users.user_id, users.user_name, users.account_icon from users as users, hot_switch_relations as switch where users.user_id=switch.to_user and switch.from_user=$1`
	rows, err := u.DB.Query(context.Background(), query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []*entitie.UserSwitcher
	for rows.Next() {
		item := &entitie.UserSwitcher{}
		err := rows.Scan(&item.Id, &item.UserName, &item.Icon)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

func (u *Service) CanUserSwitchToId(from, to int) bool {
	r := entitie.HotSwitchRelations{}
	query := `select * from hot_switch_relations where from_user=$1 and to_user=$2`
	row := u.DB.QueryRow(context.Background(), query, from, to)

	fmt.Println(query)
	fmt.Println(from)
	fmt.Println(to)

	err := row.Scan(&r.FromUser, &r.ToUser)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
