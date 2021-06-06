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
		logger.Debug("Robin> " + r.Text)
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}
	}
}

func main() {
	msg := make(chan types.Message)
	reply := make(chan types.Reply)

	hub := network.Hub{}
	hub.Init()

	ai := ai.AICore{}

	robin_server := new(network.RobibWSServer)
	robin_server.Init(logger, hub)

	wg := new(sync.WaitGroup)
	wg.Add(3)

	go robin_server.Start(wg, msg)
	go ai.Start(wg, msg, reply, logger)
	go executor(wg, reply)
	wg.Wait()

	logger.Infoln("main finished")
}
