package network

import (
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

func (h *Hub) ProcessWSSignal() {

}

func (h *Hub) Disconnect(ip string) {

}
