package api

import (
	"fleetpilot/common/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func WsHandler(ctx *gin.Context) {
	//处理ws消息程序

	// 升级http协议到ws
	c, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Error("ws协议升级失败", err)
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{
			"meg": err,
		})
	}

	defer c.Close()

	// 返回正常消息
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			logger.Debug("read:", err)
			break
		}
		logger.Debug("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			logger.Error("write:", err)
			break
		}
	}
}
