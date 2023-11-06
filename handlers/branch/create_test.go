package branch

import (
	"github.com/dmidokov/rv2/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"testing"
)

func TestService_Create(t *testing.T) {
	type fields struct {
		Logger *logrus.Logger
		DB     *pgxpool.Pool
		Config *config.Configuration
	}
	type args struct {
		branchCreator branchCreator
		userProvider  userProvider
	}
	var tests []struct {
		name   string
		fields fields
		args   args
		want   http.HandlerFunc
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Logger: tt.fields.Logger,
				DB:     tt.fields.DB,
				Config: tt.fields.Config,
			}
			if got := s.Create(tt.args.branchCreator, tt.args.userProvider); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		Logger *logrus.Logger
		DB     *pgxpool.Pool
		Config *config.Configuration
	}
	var tests []struct {
		name string
		args args
		want Service
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.Logger, tt.args.DB, tt.args.Config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_DeleteBranch(t *testing.T) {
	type fields struct {
		Logger *logrus.Logger
		DB     *pgxpool.Pool
		Config *config.Configuration
	}
	type args struct {
		branchProvider DeleteProvider
		userProvider   userProvider
	}
	var tests []struct {
		name   string
		fields fields
		args   args
		want   http.HandlerFunc
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Logger: tt.fields.Logger,
				DB:     tt.fields.DB,
				Config: tt.fields.Config,
			}
			if got := s.DeleteBranch(tt.args.branchProvider, tt.args.userProvider); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteBranch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Get(t *testing.T) {
	type fields struct {
		Logger *logrus.Logger
		DB     *pgxpool.Pool
		Config *config.Configuration
	}
	type args struct {
		branchGetter branchGetter
		userProvider userProvider
	}
	var tests []struct {
		name   string
		fields fields
		args   args
		want   http.HandlerFunc
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Logger: tt.fields.Logger,
				DB:     tt.fields.DB,
				Config: tt.fields.Config,
			}
			if got := s.Get(tt.args.branchGetter, tt.args.userProvider); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
