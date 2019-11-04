package imil

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	wsConn    *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan byte
	mutex     sync.Mutex
	isClosed  bool
}

func InitConnection(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConn:    wsConn,
		inChan:    make(chan []byte, 1000),
		outChan:   make(chan []byte, 1000),
		closeChan: make(chan byte, 1),
	}

	go conn.ReadLoop()
	go conn.WriteLoop()

	return
}

//API
//读取消息
func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

//写消息
func (conn *Connection) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

//关闭websocket
func (conn *Connection) Close() {
	conn.wsConn.Close()
	conn.mutex.Lock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
	}

}

//实现读
func (conn *Connection) ReadLoop() {
	var (
		data []byte
		err  error
	)

	for {
		if _, data, err = conn.wsConn.ReadMessage(); err != nil {
			conn.Close()
		}
		select {
		case conn.inChan <- data:
		case <-conn.closeChan:
			conn.Close()
		}
	}
}

//实现写
func (conn *Connection) WriteLoop() {
	var (
		data []byte
		err  error
	)
	for {
		data = <-conn.outChan
		if err = conn.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
			conn.Close()
		}
	}
}
