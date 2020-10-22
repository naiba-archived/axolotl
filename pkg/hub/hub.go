package hub

import (
	"encoding/json"
	"sync"

	"github.com/gofiber/websocket/v2"

	"github.com/naiba/axolotl/internal/model"
)

type Topic struct {
	Conns map[string]*websocket.Conn
	Lang  string
	Lock  sync.Mutex
}

type TopicSerialize struct {
	Lang string   `json:"lang,omitempty"`
	User []string `json:"user,omitempty"`
}

type Message struct {
	Data  interface{}
	Topic string
	From  string
}

type Subscription struct {
	Conn  *websocket.Conn
	User  string
	Topic string
}

type Hub struct {
	TopicsLock  sync.RWMutex
	Topics      map[string]*Topic
	Broadcast   chan Message
	Subscribe   chan Subscription
	UnSubscribe chan Subscription
}

func New() *Hub {
	return &Hub{
		Broadcast:   make(chan Message),
		Subscribe:   make(chan Subscription),
		UnSubscribe: make(chan Subscription),
		Topics:      make(map[string]*Topic),
	}
}

func (h *Hub) SendMsgTo(topicID, user string, msgType int, data []byte) {
	h.TopicsLock.RLock()
	defer h.TopicsLock.RUnlock()

	if topic, has := h.Topics[topicID]; has {
		topic.Lock.Lock()
		defer topic.Lock.Unlock()
		for onlineUser, conn := range topic.Conns {
			if user == onlineUser {
				if conn != nil && conn.Conn != nil {
					conn.WriteMessage(msgType, data)
				}
				return
			}
		}
	}
}

func (h *Hub) HasUser(topicID string, user string) bool {
	h.TopicsLock.RLock()
	defer h.TopicsLock.RUnlock()

	if topic, has := h.Topics[topicID]; has {
		topic.Lock.Lock()
		defer topic.Lock.Unlock()
		for onlineUser := range topic.Conns {
			if user == onlineUser {
				return true
			}
		}
	}

	return false
}

func (h *Hub) Serialize(topicID string, skip string) TopicSerialize {
	h.TopicsLock.RLock()
	defer h.TopicsLock.RUnlock()
	var topicSerialize TopicSerialize
	if topic, has := h.Topics[topicID]; has {
		topic.Lock.Lock()
		defer topic.Lock.Unlock()
		topicSerialize.Lang = topic.Lang
		for user := range topic.Conns {
			if user != skip {
				topicSerialize.User = append(topicSerialize.User, user)
			}
		}
	}
	return topicSerialize
}

func (h *Hub) UpdateLang(topicID string, lang string) {
	h.TopicsLock.Lock()
	defer h.TopicsLock.Unlock()
	h.Topics[topicID].Lang = lang
}

func (h *Hub) Serve() {
	for {
		select {
		case s := <-h.Subscribe:
			h.TopicsLock.Lock()
			topic := h.Topics[s.Topic]
			h.TopicsLock.Unlock()

			if topic == nil {
				topic = &Topic{
					Conns: make(map[string]*websocket.Conn),
				}
				topic.Lock.Lock()
				h.Topics[s.Topic] = topic
			} else {
				topic.Lock.Lock()
			}
			h.Topics[s.Topic].Conns[s.User] = s.Conn
			topic.Lock.Unlock()

		case s := <-h.UnSubscribe:
			h.TopicsLock.Lock()
			topic := h.Topics[s.Topic]
			h.TopicsLock.Unlock()

			topic.Lock.Lock()
			if topic != nil {
				if _, ok := topic.Conns[s.User]; ok {
					delete(topic.Conns, s.User)
					if s.Conn != nil && s.Conn.Conn != nil {
						s.Conn.Conn.Close()
					}
					if len(topic.Conns) == 0 {
						h.TopicsLock.Lock()
						delete(h.Topics, s.Topic)
						h.TopicsLock.Unlock()
					}
				}
			}
			topic.Lock.Unlock()

		case m := <-h.Broadcast:
			h.TopicsLock.RLock()
			topic := h.Topics[m.Topic]
			h.TopicsLock.RUnlock()

			topic.Lock.Lock()
			for u, c := range topic.Conns {
				if u == m.From {
					continue
				}
				var err error
				if c != nil && c.Conn != nil {
					if data, ok := m.Data.(model.WsMsg); ok {
						data.From = m.From
						content, err := json.Marshal(data)
						if err == nil {
							err = c.WriteMessage(websocket.TextMessage, content)
						}
					} else {
						if data, ok := m.Data.([]byte); ok {
							err = c.WriteMessage(websocket.BinaryMessage, data)
						} else {
							err = c.WriteMessage(websocket.BinaryMessage, []byte{})
						}
					}
					if err == nil {
						continue
					}
				}
				delete(topic.Conns, u)
				if len(topic.Conns) == 0 {
					h.TopicsLock.Lock()
					delete(h.Topics, m.Topic)
					h.TopicsLock.Unlock()
				}
			}
			topic.Lock.Unlock()
		}
	}
}
