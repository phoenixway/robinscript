package types

import "github.com/gorilla/websocket"

type Message struct {
	Text string
	Ws   *websocket.Conn
}

type Reply struct {
	Text string
	Ws   *websocket.Conn
}
