package handlers

import (
	"errors"
	"net/http"
)

func (hm *Service) GetSessionUserId(r *http.Request) (int, error) {

	if auth, ok := hm.CookieStore.Get(r, "authenticated"); ok && auth.(bool) {

		if userId, ok := hm.CookieStore.Get(r, "userid"); ok {
			return userId.(int), nil
		}

		return -1, errors.New("user id is empty")
	}
	return -1, errors.New("user not authorized")
}

func (hm *Service) GetSessionOrganizationId(r *http.Request) (int, error) {

	if auth, ok := hm.CookieStore.Get(r, "authenticated"); ok && auth.(bool) {

		if orgId, ok := hm.CookieStore.Get(r, "organizationid"); ok {
			return orgId.(int), nil
		}

		return -1, errors.New("organization id is empty")
	}
	return -1, errors.New("user not authorized")
}
