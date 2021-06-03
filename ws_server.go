package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/phoenixway/robinscript/users"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type NetworkHub struct {
	UsernameByIP map[string]string
}

func (h *NetworkHub) Init() {
	h.UsernameByIP = map[string]string{}
}

func (h *NetworkHub) IsNewConnection(ip string) bool {
	if h.UsernameByIP[ip] != "" {
		return false
	} else {
		return true
	}
}

func (h *NetworkHub) IsAuthenticated(ip string) bool {
	return (h.UsernameByIP[ip] != "") && (h.UsernameByIP[ip] != "Guest")
}

func (h *NetworkHub) DoLogin(ip, login, pass string) users.UserAccount {

}

func (h *NetworkHub) ProcessWSSignal() {

}

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

func handleText(ws *websocket.Conn, text string) {
	answer := fmt.Sprintf("You said: %s", text)
	err := ws.WriteMessage(websocket.TextMessage, []byte(answer))
	logger.Debug("Server> " + answer)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

func handleWebsockets(c echo.Context) error {

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
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
			handleText(ws, string(msg))
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
