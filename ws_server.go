package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/phoenixway/robinscript/aicore"
	"github.com/phoenixway/robinscript/network"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	upgrader = websocket.Upgrader{}
	logger   = &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &prefixed.TextFormatter{
			ForceColors:     true,
			TimestampFormat: "15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		},
	}
)

func authenticate(login string) {
	logger.Debug("Autothentication started!")
	logger.Debug(fmt.Sprintf("User's login: %s", login))
}

func handleCommand(com string) {
	logger.Debug("Command received: " + com)
	isAuth := regexp.MustCompile(`/auth (\w+)`)
	switch {
	case isAuth.MatchString(com):
		authenticate(isAuth.FindStringSubmatch(com)[1])
	default:
		return
	}
}

func handleText(ws *websocket.Conn, ip, text string) {

	var hub = network.Hub{}
	hub.Init()
	u := hub.UserByIP(ip)
	//TODO: change it to sending to channel
	answer := aicore.ProcessMessage(u, text)
	answer = fmt.Sprintf("Client said: %s", text)
	err := ws.WriteMessage(websocket.TextMessage, []byte(answer))
	logger.Debug("Server> " + answer)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

func handleWebsockets(c echo.Context) error {

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	ip, _, err := net.SplitHostPort(c.Request().RemoteAddr)

	if err != nil {
		logger.Error(err)
		return err
	}
	defer ws.Close()

	logger.Debug("Client connected!")
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}
		logger.Debug("Client> " + string(msg))

		isCommand := regexp.MustCompile(`^/.+`)
		switch {
		case isCommand.MatchString(string(msg)):
			handleCommand(string(msg))
		default:
			handleText(ws, ip, string(msg))
		}

	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Static("/", "../public")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	e.GET("/ws", handleWebsockets)
	e.Logger.Fatal(e.Start(":1323"))
}
