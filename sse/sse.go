package sse

import (
	"fmt"
	"net/http"
)

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

func (s *EventService) SendToQueue() {

}
