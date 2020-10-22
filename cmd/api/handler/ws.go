package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	"github.com/naiba/axolotl/internal/model"
	"github.com/naiba/axolotl/pkg/hub"
	"github.com/naiba/axolotl/pkg/util"
)

func NotInRoom(pubsub *hub.Hub) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roomID := c.Params("conferenceID")
		user := c.Locals(model.KeyAuthorizedUser).(model.User).Nickname

		// 判断用户是否在线
		if pubsub.HasUser(roomID, user) {
			c.WriteString("User alreay in room")
			c.Status(http.StatusBadRequest)
			return nil
		}

		return c.Next()
	}
}

func WS(pubsub *hub.Hub) func(c *websocket.Conn) {
	return func(c *websocket.Conn) {
		roomID := c.Params("conferenceID")
		user := c.Locals(model.KeyAuthorizedUser).(model.User).Nickname

		pubsub.Subscribe <- hub.Subscription{
			Conn:  c,
			User:  user,
			Topic: roomID,
		}

		// 新加入用户：发送当前会议室的最新信息
		topic := pubsub.Serialize(roomID, user)
		if topic.Lang == "" {
			topic.Lang = "go"
		}
		data, err := json.Marshal(model.WsMsg{
			Type: model.MsgTypeInitRoom,
			Data: topic,
			From: user,
		})
		if err == nil {
			c.WriteMessage(websocket.TextMessage, data)
		}

		// 其他已加入会议用户：创建 Peer 的信令
		pubsub.Broadcast <- hub.Message{
			Topic: roomID,
			Data: model.WsMsg{
				Type: model.MsgTypeCreatePeer,
			},
			From: user,
		}

		var m model.WsMsg
		for {
			mType, data, err := c.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err) && !websocket.IsUnexpectedCloseError(err) {
					util.Infof(0, "websocket error: %+v", err)
				}
				pubsub.UnSubscribe <- hub.Subscription{
					Conn:  c,
					User:  user,
					Topic: roomID,
				}
				break
			}
			if mType == websocket.TextMessage {
				if err = json.Unmarshal(data, &m); err == nil {
					// 指哪打哪的 WebRTC 信令
					if m.Type == model.MsgTypePeer {
						m.From = user
						to := m.To
						m.To = ""
						data, err = json.Marshal(m)
						if err != nil {
							continue
						}
						pubsub.SendMsgTo(roomID, to, websocket.TextMessage, data)
						continue
					}
					// 其他消息
					if m.Type == model.MsgTypeChooseProgrammingLanguage {
						pubsub.UpdateLang(roomID, m.Data.(string))
					}
					pubsub.Broadcast <- hub.Message{
						Topic: roomID,
						Data:  m,
						From:  user,
					}
				}
			} else if mType == websocket.BinaryMessage {
				pubsub.Broadcast <- hub.Message{
					Topic: roomID,
					Data:  data,
					From:  user,
				}
			}
		}
	}
}
