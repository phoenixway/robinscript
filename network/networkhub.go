package network

import (
	"github.com/phoenixway/robinscript/users"
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
	return users.UserAccount{}
}

func (h *NetworkHub) ProcessWSSignal() {

}

func (h *NetworkHub) Disconnect(ip string) {

}
