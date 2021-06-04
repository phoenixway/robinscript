package network

import (
	"testing"
)

func TestHub(t *testing.T) {
	var h = Hub{}
	h.Init()
	ip := "192.168.1.1"
	message := "hello"
	h.ProcessWSSignal(ip, message)
	/*
		connect new user
		check if it added as guest
		login
		check
		send signal
		check
		disconnect
		check
	*/
}
