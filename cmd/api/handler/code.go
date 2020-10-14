package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/naiba/helloengineer/internal/model"
	"github.com/naiba/helloengineer/pkg/hub"
)

func RunCode(conf *model.Config, h *hub.Hub) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req model.RunCodeRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}
		body, err := json.Marshal(req)
		if err != nil {
			return err
		}

		client := &http.Client{}
		remoteReq, err := http.NewRequest("POST", conf.Code.Endpoint+"/api/run", bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		remoteReq.Header.Add("Authorization", "Basic "+conf.Code.Authorization)
		remoteReq.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(remoteReq)
		if err != nil {
			return err
		}
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		data, err := json.Marshal(model.WsMsg{
			Type: model.MsgTypeExecResult,
			Data: string(body),
		})
		if err != nil {
			return err
		}
		h.Broadcast <- hub.Message{
			Topic: req.Room,
			Data:  data,
			From:  nil,
		}
		return nil
	}
}

func ListRunner(conf *model.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		client := &http.Client{}
		remoteReq, err := http.NewRequest("GET", conf.Code.Endpoint+"/api/list", nil)
		if err != nil {
			return err
		}
		remoteReq.Header.Add("Authorization", "Basic "+conf.Code.Authorization)

		resp, err := client.Do(remoteReq)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return c.JSON(model.Response{
			Data: string(body),
		})
	}
}
