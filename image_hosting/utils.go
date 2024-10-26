package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

func LskyRefreshToken() (string, error) {
	Token.Lock()
	defer Token.Unlock()

	agent := LskyBaseAgent(fiber.AcquireAgent(), fiber.MethodPost, "/tokens")
	defer fiber.ReleaseAgent(agent)

	agent.JSON(fiber.Map{"email": Config.ProxyEmail, "password": Config.ProxyPassword})

	if err := agent.Parse(); err != nil {
		return "", err
	}

	code, body, errs := agent.Bytes()
	if len(errs) != 0 {
		return "", errs[0]
	}
	if code != 200 {
		message := fmt.Sprintf(`{"code": %v}`, code)
		return "", fiber.NewError(fiber.StatusInternalServerError, message)
	}

	var lskyToken LskyToken
	err := json.Unmarshal(body, &lskyToken)
	if err != nil {
		return "", err
	}
	Token.data = lskyToken.Data.Token
	return lskyToken.Data.Token, nil
}

func LskyBaseAgent(agent *fiber.Agent, method string, path string) *fiber.Agent {
	agent.Set("Accept", fiber.MIMEApplicationJSON).
		UserAgent("fiber").
		Timeout(time.Second * 10).Reuse()
	req := agent.Request()
	req.Header.SetMethod(method)
	req.SetRequestURI(Config.ProxyUrl + path)
	return agent
}
