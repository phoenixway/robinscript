package network

import (
	"github.com/phoenixway/robinscript/ai"
	"github.com/phoenixway/robinscript/users"
)

type ConnectedUsers struct {
	UsernameByIP map[string]string
}

func (h *ConnectedUsers) Init() {
	h.UsernameByIP = map[string]string{}
}

func (h *ConnectedUsers) IsNewConnection(ip string) bool {
	if h.UsernameByIP[ip] != "" {
		return false
	} else {
		return true
	}
}

func (h *ConnectedUsers) IsAuthenticated(ip string) bool {
	return (h.UsernameByIP[ip] != "") && (h.UsernameByIP[ip] != "Guest")
}

func (h *ConnectedUsers) DoLogin(ip, login, pass string) users.UserAccount {
	return users.UserAccount{}
}

func (h *ConnectedUsers) ProcessWSSignal(ip, message string) {
	//TODO: remove this func
	u := h.UserByIP(ip)
	//TODO: change it to sending to channel
	ai.ProcessMessage(u, message)
}

func (h *ConnectedUsers) Disconnect(ip string) {

}

func (h *ConnectedUsers) UserByIP(ip string) *users.UserAccount {
	if h.IsNewConnection(ip) {
		//and if it's not a command to login and there is no guest account - create guest account
		//if its command to login restore or create user account
	}
	return nil
}
