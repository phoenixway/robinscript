package network

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/phoenixway/robinscript/types"
	"github.com/sirupsen/logrus"
)

type RobibWSServer struct {
	upgrader websocket.Upgrader
	logger   *logrus.Logger
	msg      chan types.Message
	wg       *sync.WaitGroup
}

func (r *RobibWSServer) Init(logger *logrus.Logger, hub Hub) {
	r.logger = logger
}

func (r *RobibWSServer) Start(wg *sync.WaitGroup, msg chan types.Message) {
	defer wg.Done()
	r.logger.Debugf("Websocket server started!")
	r.msg = msg
	r.wg = wg
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

	r.logger.Info(fmt.Sprintf("User from %s is connected!", ip))
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			r.logger.Error(err)
			os.Exit(1)
		}
		m := types.Message{string(msg), ws}
		r.logger.Info("User> " + string(msg))

		isCommand := regexp.MustCompile(`^/.+`)
		switch {
		case isCommand.MatchString(string(msg)):
			r.handleCommand(string(msg))
		default:
			r.handleText(m)
		}

	}
}
func (r *RobibWSServer) handleCommand(s string) {
	r.logger.Debug("Command received: " + s)
	isAuth := regexp.MustCompile(`/auth (\w+)`)
	switch {
	case isAuth.MatchString(s):
		//authenticate(isAuth.FindStringSubmatch(com)[1])
	default:
		//TODO: add /quit command
		close(r.msg)
		r.logger.Debugf("Websocket server stopped!")
		return
	}

}

func (r *RobibWSServer) handleText(m types.Message) {

	var hub = Hub{}
	hub.Init()
	//u := hub.UserByIP(ip)
	//TODO: change it to sending to channel
	// answer := ai.ProcessMessage(u, msg)
	// answer = fmt.Sprintf("Client said: %s", msg)
	r.msg <- m
	//err := ws.WriteMessage(websocket.TextMessage, []byte(answer))
	// r.logger.Debug("Server> " + answer)
	// if err != nil {
	// 	r.logger.Error(err)
	// 	os.Exit(1)
	// }
}
