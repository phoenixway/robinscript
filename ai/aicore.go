package ai

import (
	"fmt"
	"os"
	"sync"

	"github.com/phoenixway/robinscript/types"
	"github.com/phoenixway/robinscript/users"
	"github.com/sirupsen/logrus"
)

type AICore struct {
	logger *logrus.Logger
	msg    chan types.Message
	reply  chan types.Reply
	wg     *sync.WaitGroup
}

func (ai *AICore) Start(wg *sync.WaitGroup, msg chan types.Message, reply chan types.Reply, logger *logrus.Logger) {
	ai.msg = msg
	ai.reply = reply
	ai.logger = logger
	ai.logger.Debugf("AI started!")
	defer wg.Done()
	var ok bool
	var m types.Message
	for {
		m, ok = <-msg
		if !ok {
			wg.Done()
			logger.Debugln("AI is finished!")
			os.Exit(0)
			return
		}
		logger.Debugf("AI got '%s'!", m.Text)
		answer := fmt.Sprintf("echo reply from '%s'", m.Text)
		r := types.Reply{answer, m.Ws}
		reply <- r
	}

}

func ProcessMessage(u *users.UserAccount, m string) string {
	return ""
}
