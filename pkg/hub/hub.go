package hub

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Message struct {
	Type  uint
	Data  []byte
	Topic string
	From  *websocket.Conn
}

type Subscription struct {
	Conn  *websocket.Conn
	Topic string
}

type Hub struct {
	TopicsLock  sync.RWMutex
	Topics      map[string]map[*websocket.Conn]bool
	Broadcast   chan Message
	Subscribe   chan Subscription
	UnSubscribe chan Subscription
}

func New() *Hub {
	return &Hub{
		Broadcast:   make(chan Message),
		Subscribe:   make(chan Subscription),
		UnSubscribe: make(chan Subscription),
		Topics:      make(map[string]map[*websocket.Conn]bool),
	}
}

func (h *Hub) Serve() {
	for {
		select {
		case s := <-h.Subscribe:
			h.TopicsLock.Lock()
			connections := h.Topics[s.Topic]

			if connections == nil {
				connections = make(map[*websocket.Conn]bool)
				h.Topics[s.Topic] = connections
			}

			h.Topics[s.Topic][s.Conn] = true
			h.TopicsLock.Unlock()
		case s := <-h.UnSubscribe:
			h.TopicsLock.Lock()
			connections := h.Topics[s.Topic]

			if connections != nil {
				if _, ok := connections[s.Conn]; ok {
					delete(connections, s.Conn)
					s.Conn.Conn.Close()

					if len(connections) == 0 {
						delete(h.Topics, s.Topic)
					}
				}
			}
			h.TopicsLock.Unlock()
		case m := <-h.Broadcast:
			h.TopicsLock.RLock()
			connections := h.Topics[m.Topic]
			h.TopicsLock.RUnlock()

			for c := range connections {
				if c.Conn != nil {
					if c == m.From {
						continue
					}
					var err error
					switch m.Type {
					case websocket.BinaryMessage:
						err = c.WriteMessage(websocket.BinaryMessage, m.Data)
					default:
						err = c.WriteMessage(websocket.TextMessage, m.Data)
					}
					if err == nil {
						continue
					}
				}
				h.TopicsLock.Lock()
				delete(connections, c)
				if len(connections) == 0 {
					delete(h.Topics, m.Topic)
				}
				h.TopicsLock.Unlock()
			}
		}
	}
}
