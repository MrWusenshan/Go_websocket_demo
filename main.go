package main

import (
	"github.com/gorilla/websocket"
	"go-websocket/imil"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		//允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(writer http.ResponseWriter, request *http.Request) {
	var (
		err    error
		wsConn *websocket.Conn
		conn   *imil.Connection
		data []byte
	)

	if wsConn, err = upgrader.Upgrade(writer, request, nil); err != nil {
		return
	}

	conn, err = imil.InitConnection(wsConn)
	if err != nil {
		goto ERR
	}

	go func() {
		if err := conn.WriteMessage([]byte("heartbeat")); err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}()

	for {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		if err = conn.WriteMessage(data); err != nil {
			goto ERR
		}
	}
ERR:
	//TODO:关闭连接
	conn.Close()
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe("0.0.0.0:7777", nil)
}
