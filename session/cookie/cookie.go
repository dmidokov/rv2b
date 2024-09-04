package cookie

import (
	"github.com/gorilla/sessions"
	"net/http"
)

const (
	Authenticated  = "authenticated"
	UserId         = "userId"
	OrganizationId = "organizationId"
	SwitchedTo     = "switchedTo"
	SelectedBranch = "selectedBranch"
)

type Service struct {
	CookieStore *sessions.CookieStore
	Secret      string
}

func New(secret string) *Service {
	return &Service{
		CookieStore: sessions.NewCookieStore([]byte(secret)),
		Secret:      secret,
	}
}

func (s *Service) SetMaxAge(maxAge int) {
	s.CookieStore.Options.MaxAge = maxAge
}

func (s *Service) Save(r *http.Request, w http.ResponseWriter, data map[string]interface{}) bool {
	session, _ := s.CookieStore.Get(r, s.Secret)

	for k, v := range data {
		session.Values[k] = v
	}
	err := session.Save(r, w)

	if err != nil {
		return false
	}

	return true
}

func (s *Service) GetByKey(r *http.Request, key string) (interface{}, bool) {
	session, err := s.CookieStore.Get(r, s.Secret)

	if err != nil {
		return "", false
	}

	if v, ok := session.Values[key]; ok {
		return v, true
	}

	return nil, false
}

func (s *Service) Get(r *http.Request) (map[interface{}]interface{}, error) {
	session, err := s.CookieStore.Get(r, s.Secret)
	if err != nil {
		return nil, err
	}

	return session.Values, nil
}
