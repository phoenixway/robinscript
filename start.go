package main

import (
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/phoenixway/robinscript/ai"
	"github.com/phoenixway/robinscript/network"
	"github.com/phoenixway/robinscript/types"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	logger = &logrus.Logger{
		Out:          os.Stderr,
		Hooks:        map[logrus.Level][]logrus.Hook{},
		Formatter:    &prefixed.TextFormatter{ForceColors: true, TimestampFormat: "15:04:05", FullTimestamp: true, ForceFormatting: true},
		ReportCaller: false,
		Level:        logrus.DebugLevel,
		ExitFunc: func(int) {
		},
	}
)

func executor(wg *sync.WaitGroup, reply chan types.Reply) {
	logger.Debug("Executor started!")
	defer wg.Done()
	for {
		r := <-reply
		logger.Debugf("Executor got \"%s\"!", r.Text)
		err := r.Ws.WriteMessage(websocket.TextMessage, []byte(r.Text))
		logger.Info("Robin> " + r.Text)
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}
	}
}

func main() {
	messages := make(chan types.Message)
	replies := make(chan types.Reply)

	connectedUsers := network.ConnectedUsers{}
	connectedUsers.Init()

	ai := ai.AICore{}

	ws_server := new(network.RobibWSServer)
	ws_server.Init(logger, connectedUsers)

	wg := new(sync.WaitGroup)
	wg.Add(3)

	go ws_server.Start(wg, messages)
	go ai.Start(wg, messages, replies, logger)
	go executor(wg, replies)
	wg.Wait()

	logger.Infoln("main() is finished!")
}
