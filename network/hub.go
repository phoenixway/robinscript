package network

import (
	"github.com/phoenixway/robinscript/aicore"
	"github.com/phoenixway/robinscript/users"
)

type Hub struct {
	UsernameByIP map[string]string
}

func (h *Hub) Init() {
	h.UsernameByIP = map[string]string{}
}

func (h *Hub) IsNewConnection(ip string) bool {
	if h.UsernameByIP[ip] != "" {
		return false
	} else {
		return true
	}
}

func (h *Hub) IsAuthenticated(ip string) bool {
	return (h.UsernameByIP[ip] != "") && (h.UsernameByIP[ip] != "Guest")
}

func (h *Hub) DoLogin(ip, login, pass string) users.UserAccount {
	return users.UserAccount{}
}

func (h *Hub) ProcessWSSignal(ip, message string) {
	//TODO: remove this func
	u := h.UserByIP(ip)
	//TODO: change it to sending to channel
	aicore.ProcessMessage(u, message)
}

func (h *Hub) Disconnect(ip string) {

}

func (h *Hub) UserByIP(ip string) *users.UserAccount {
	if h.IsNewConnection(ip) {
		//and if it's not a command to login and there is no guest account - create guest account
		//if its command to login restore or create user account
	}
	return nil
}
