package branch

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

type SessionStorage interface {
	Save(r *http.Request, w http.ResponseWriter, data map[string]interface{}) bool
}

func New(DB *pgxpool.Pool, CookieStore SessionStorage, Log *logrus.Logger) *Service {
	return &Service{
		DB:          DB,
		CookieStore: CookieStore,
		Log:         Log,
	}
}

func (o *Service) GetAll(userId int) ([]*e.Branch, error) {
	query := `select * from branches where branch_id in (
    select branch_id from user_branches where user_id=$1
    );`
	rows, err := o.DB.Query(context.Background(), query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRows(rows)
}

func scanRows(rows pgx.Rows) ([]*e.Branch, error) {
	var result []*e.Branch
	for rows.Next() {
		item := &e.Branch{}
		err := rows.Scan(&item.Id, &item.OrgId, &item.Name, &item.Address, &item.Phone, &item.WorkTime, &item.CreateTime, &item.UpdateTime)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

func (o *Service) Create(branch *e.Branch, userId int) (*e.Branch, error) {
	query := `
		INSERT INTO branches
			(branch_name, organization_id, address, phone, work_time, create_time, update_time) 
		VALUES
			($1, $2, $3, $4, $5, $6, $6);`

	tag, err := o.DB.Exec(context.Background(), query, branch.Name, branch.OrgId, branch.Address, branch.Phone, branch.WorkTime, time.Now().Unix())
	if err != nil {
		return nil, err
	}
	logrus.Info(tag.RowsAffected())

	newBranch, err := o.GetByNameAndOrgId(branch.Name, branch.OrgId)
	if err != nil {
		return nil, err
	}

	query = `
		INSERT INTO user_branches (user_id, branch_id)
	    VALUES ($1, $2);`

	tag, err = o.DB.Exec(context.Background(), query, userId, newBranch.Id)
	if err != nil {
		return nil, err
	}

	return newBranch, nil
}

func (o *Service) GetByNameAndOrgId(name string, orgId int) (*e.Branch, error) {
	query := `select * from branches where branch_name=$1 and organization_id=$2`
	row := o.DB.QueryRow(context.Background(), query, name, orgId)

	branch := &e.Branch{}

	err := row.Scan(
		&branch.Id,
		&branch.OrgId,
		&branch.Name,
		&branch.Address,
		&branch.Phone,
		&branch.WorkTime,
		&branch.CreateTime,
		&branch.UpdateTime,
	)

	if err != nil {
		return nil, err
	}

	return branch, nil
}

func (o *Service) Delete(branchId int, orgId int) error {
	query := `
		DELETE FROM branches 
		WHERE branch_id=$1 and organization_id=$2`

	tag, err := o.DB.Exec(context.Background(), query, branchId, orgId)
	if err != nil {
		return err
	}
	logrus.Info("Branches deleted: ", tag.RowsAffected())

	return nil
}

func (o *Service) SetActive() error {

	return nil
}
