package interfaces

import "github.com/gorilla/websocket"

// ToolHandler 是所有工具必须实现的接口
type ToolHandler interface {
	GetToolName() string
	Executed(conn *websocket.Conn, msg []byte) error
}
