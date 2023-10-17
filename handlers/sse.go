package handlers

import (
	"github.com/dmidokov/rv2/response"
	"github.com/dmidokov/rv2/sse"
	"net/http"
	"path"
)

func (hm *Service) sseHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		log := hm.Logger
		responses := resp.Service{Writer: &w, Logger: log, Operation: "sse.sseHandler"}

		if auth, ok := hm.CookieStore.Get(r, "authenticated"); !ok || !auth.(bool) {
			log.Warning("User is not authorized")
			responses.Unauthorized()
			return
		}

		if _, ok := hm.CookieStore.Get(r, "userid"); !ok {
			log.Warning("User is not authorized")
			responses.Unauthorized()
			return
		}

		eventName := sse.EventName(path.Base(r.URL.Path))

		v, _ := hm.CookieStore.Get(r, "userid")
		userId := v.(int)

		flusher, _ := w.(http.Flusher)
		a := make(map[int]sse.Receiver)
		if v, ok := hm.SSE.Receivers[eventName]; ok {
			a = v
		}
		a[userId] = sse.Receiver{Wr: &w, Fl: flusher}
		hm.SSE.Receivers[eventName] = a
		flusher.Flush()

		for {
			select {
			case <-r.Context().Done():
				delete(hm.SSE.Receivers[eventName], userId)
				return
			}
		}

	}
}
