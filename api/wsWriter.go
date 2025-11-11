package api

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WsWriter interface {
	Write([]byte) error
	WriteJSON(interface{}) error
}

type wsWriter struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func NewWsWriter(conn *websocket.Conn) WsWriter {
	return &wsWriter{conn: conn}
}

func (w *wsWriter) Write(b []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.conn.WriteMessage(websocket.TextMessage, b)
}

func (w *wsWriter) WriteJSON(v interface{}) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.conn.WriteJSON(v)
}
