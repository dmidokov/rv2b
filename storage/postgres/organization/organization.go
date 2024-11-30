package organization

import (
	"context"
	e "github.com/dmidokov/rv2/lib/entitie"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
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

func (o *Service) GetByHostName(hostName string) (*e.Organization, error) {
	hostWithoutPort := strings.Split(hostName, ":")[0]
	query := "select * from organizations where host=$1"
	row := o.DB.QueryRow(context.Background(), query, hostWithoutPort)

	organization := &e.Organization{}

	err := row.Scan(
		&organization.Id,
		&organization.Name,
		&organization.Host,
		&organization.CreateTime,
		&organization.UpdateTime,
		&organization.Creator,
	)

	if err != nil {
		return nil, err
	}

	return organization, nil
}

func (o *Service) GetAll() ([]*e.Organization, error) {
	query := `select * from organizations`
	rows, err := o.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRows(rows)
}

func (o *Service) Create(org *e.Organization) (*e.Organization, error) {
	query := `
		INSERT INTO organizations
			(organization_name, host, create_time, update_time, creator) 
		VALUES
			($1, $2, $3, $3, $4);`

	tag, err := o.DB.Exec(context.Background(), query, org.Name, org.Host, time.Now().Unix(), org.Creator)
	if err != nil {
		return nil, err
	}
	logrus.Info(tag.RowsAffected())

	organization, err := o.GetByHostName(org.Host)
	if err != nil {
		return nil, err
	}

	return organization, nil
}

func scanRows(rows pgx.Rows) ([]*e.Organization, error) {
	var result []*e.Organization
	for rows.Next() {
		item := &e.Organization{}
		err := rows.Scan(&item.Id, &item.Name, &item.Host, &item.CreateTime, &item.UpdateTime, &item.Creator)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

func (o *Service) Delete(orgId int) error {
	query := `
		DELETE FROM organizations 
		WHERE organization_id=$1`

	tag, err := o.DB.Exec(context.Background(), query, orgId)
	if err != nil {
		return err
	}
	logrus.Info("Organizations deleted: ", tag.RowsAffected())

	return nil
}

func (o *Service) GetById(id int) (*e.Organization, error) {
	query := `select * from organizations where organization_id=$1`
	row := o.DB.QueryRow(context.Background(), query, id)

	organization := &e.Organization{}

	err := row.Scan(
		&organization.Id,
		&organization.Name,
		&organization.Host,
		&organization.CreateTime,
		&organization.UpdateTime,
		&organization.Creator,
	)

	if err != nil {
		return nil, err
	}

	return organization, nil
}
