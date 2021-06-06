package main

import (
	"os"
	"sync"

	"github.com/phoenixway/robinscript/ai"
	"github.com/phoenixway/robinscript/network"
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

func main() {
	hub := network.Hub{}
	hub.Init()
	ai := ai.AICore{}
	robin_server := new(network.RobibWSServer)
	robin_server.Init(logger, hub)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go ai.Start()
	go robin_server.Start()
	wg.Wait()
}
