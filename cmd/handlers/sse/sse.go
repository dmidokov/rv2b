package sse

import (
	"fmt"
	resp "github.com/dmidokov/rv2/response"
	"github.com/dmidokov/rv2/session/cookie"
	"net/http"
	"path"
)

func (s *Service) SseHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		log := s.Logger
		responses := resp.Service{Writer: &w, Logger: log, Operation: "sse.sseHandler"}

		if auth, ok := s.CookieStore.GetByKey(r, cookie.Authenticated); !ok || !auth.(bool) {
			log.Warning("User is not authorized")
			responses.Unauthorized()
			return
		}

		var userId int
		if value, ok := s.CookieStore.GetByKey(r, cookie.UserId); !ok {
			log.Warning("User is not authorized")
			responses.Unauthorized()
			return
		} else {
			userId = value.(int)
		}

		eventName := EventName(path.Base(r.URL.Path))

		flusher, _ := w.(http.Flusher)
		a := make(map[int]Receiver)
		if v, ok := s.SSE.Receivers[eventName]; ok {
			a = v
		}
		a[userId] = Receiver{Wr: &w, Fl: flusher}
		s.SSE.Receivers[eventName] = a
		flusher.Flush()

		for {
			select {
			case <-r.Context().Done():
				delete(s.SSE.Receivers[eventName], userId)
				return
			}
		}

	}
}

type Event struct {
	Name   EventName
	Value  string
	UserId int
}

type Receiver struct {
	Wr *http.ResponseWriter
	Fl http.Flusher
}

type EventName string

type Receivers map[EventName]map[int]Receiver

type EventService struct {
	Chanel    chan Event
	Receivers Receivers
}

// Run запускает горутины для отправки сообщений из очереди подписавшимся клиентам
func (s *EventService) Run() {
	go func() {
		for {
			select {
			case message := <-s.Chanel:
				if message.UserId == 0 {
					for _, client := range s.Receivers[message.Name] {
						fmt.Fprintf(*client.Wr, "%s\n", message.Value)
						client.Fl.Flush()
					}
				} else {
					if clientsList, ok := s.Receivers[message.Name]; ok {
						if client, ok := clientsList[message.UserId]; ok {
							fmt.Fprintf(*client.Wr, "%s\n", message.Value)
							client.Fl.Flush()
						}
					}
				}
			}
		}
	}()
}

func (s *EventService) CreateEvent(name string, value string, userId int) Event {
	return Event{
		Name:   EventName(name),
		Value:  value,
		UserId: userId,
	}
}

func (s *EventService) SendToQueue(name string, value string, userId int) {
	event := s.CreateEvent(name, value, userId)
	s.Chanel <- event
}
