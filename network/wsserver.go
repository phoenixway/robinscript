package network

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/phoenixway/robinscript/ai"
	"github.com/sirupsen/logrus"
)

type RobibWSServer struct {
	upgrader websocket.Upgrader
	logger   *logrus.Logger
}

func (r *RobibWSServer) Init(logger *logrus.Logger, hub Hub) {
	r.logger = logger
}

func (r *RobibWSServer) Start() {
	r.upgrader = websocket.Upgrader{}
	e := echo.New()
	e.Use(middleware.Recover())
	e.Static("/", "../public")
	r.upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	e.GET("/ws", r.handleWebsockets)
	e.Logger.Fatal(e.Start(":1323"))

}

func (r *RobibWSServer) handleWebsockets(c echo.Context) error {

	ws, err := r.upgrader.Upgrade(c.Response(), c.Request(), nil)
	rawip := c.Request().RemoteAddr
	ip, _, err := net.SplitHostPort(rawip)

	if err != nil {
		r.logger.Error(err)
		return err
	}
	defer ws.Close()

	r.logger.Debug(fmt.Sprintf("Client from %s connected!", ip))
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			r.logger.Error(err)
			os.Exit(1)
		}
		r.logger.Debug("Client> " + string(msg))

		isCommand := regexp.MustCompile(`^/.+`)
		switch {
		case isCommand.MatchString(string(msg)):
			r.handleCommand(string(msg))
		default:
			r.handleText(ws, ip, string(msg))
		}

	}
}
func (r *RobibWSServer) handleCommand(com string) {
	r.logger.Debug("Command received: " + com)
	isAuth := regexp.MustCompile(`/auth (\w+)`)
	switch {
	case isAuth.MatchString(com):
		//authenticate(isAuth.FindStringSubmatch(com)[1])
	default:
		return
	}
}

func (r *RobibWSServer) handleText(ws *websocket.Conn, ip, msg string) {

	var hub = Hub{}
	hub.Init()
	u := hub.UserByIP(ip)
	//TODO: change it to sending to channel
	answer := ai.ProcessMessage(u, msg)
	answer = fmt.Sprintf("Client said: %s", msg)
	err := ws.WriteMessage(websocket.TextMessage, []byte(answer))
	r.logger.Debug("Server> " + answer)
	if err != nil {
		r.logger.Error(err)
		os.Exit(1)
	}
}
