package entitie

import "github.com/gorilla/sessions"

type CookieService struct {
	CookieStore *sessions.CookieStore
	Secret      string
}
