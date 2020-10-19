package handler

import (
	"encoding/json"

	"github.com/gofiber/websocket/v2"

	"github.com/naiba/helloengineer/internal/model"
	"github.com/naiba/helloengineer/pkg/hub"
	"github.com/naiba/helloengineer/pkg/util"
)

func WS(pubsub *hub.Hub) func(c *websocket.Conn) {
	return func(c *websocket.Conn) {
		roomID := c.Params("meetingID")
		pubsub.Subscribe <- hub.Subscription{
			Conn:  c,
			Topic: roomID,
		}
		var m model.WsMsg
		for {
			mType, data, err := c.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err) && !websocket.IsUnexpectedCloseError(err) {
					util.Infof(0, "websocket error: %+v", err)
				}
				c.Close()
				break
			}
			if mType == websocket.TextMessage {
				if err = json.Unmarshal(data, &m); err == nil {
					data, err := json.Marshal(m)
					if err != nil {
						util.Errorf(0, "%+v", err)
						continue
					}
					pubsub.Broadcast <- hub.Message{
						Topic: roomID,
						Data:  data,
						From:  c,
					}
				}
			} else if mType == websocket.BinaryMessage {
				pubsub.Broadcast <- hub.Message{
					Type:  websocket.BinaryMessage,
					Topic: roomID,
					Data:  data,
					From:  c,
				}
			}
		}
	}
}
