package handler

import (
	"encoding/json"

	"github.com/gofiber/websocket"

	"github.com/naiba/helloengineer/pkg/hub"
	"github.com/naiba/helloengineer/pkg/util"
)

var pubsub *hub.Hub

func init() {
	pubsub = hub.New()
	go pubsub.Serve()
}

const (
	MsgTypePeer = 0
)

type Msg struct {
	Type uint
	Data string
}

func WS() func(c *websocket.Conn) {
	return func(c *websocket.Conn) {
		c.SetPingHandler(c.PingHandler())
		c.SetPongHandler(c.PongHandler())
		roomID := c.Params("meetingID")
		pubsub.Subscribe <- hub.Subscription{
			Conn:  c,
			Topic: roomID,
		}
		var m Msg
		var err error
		for {
			err = c.ReadJSON(&m)
			if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
				c.Close()
				break
			}
			if err == nil {
				switch m.Type {
				case MsgTypePeer:
					data, err := json.Marshal(m)
					if err != nil {
						util.Errorf(0, "%+v", err)
						continue
					}
					pubsub.Broadcast <- hub.Message{
						Topic: roomID,
						Data:  data,
					}
				}
			}
		}
	}
}
