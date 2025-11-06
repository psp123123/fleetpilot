package api

import (
	"fleetpilot/common/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func WsHandler(ctx *gin.Context) {
	// 升级协议
	c, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Error("ws协议升级失败", err)
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{
			"msg": err.Error(),
		})
		return
	}
	defer c.Close()

	// 用于控制协程
	done := make(chan bool)

	// 协程1：定时发送消息（每1秒一次）
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				msg := `{"type":"chat","payload":"{\"host\":\"123\",\"scanType\":\"-sS\"}"}`
				if err := c.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
					logger.Error("write message failed:", err)
					close(done)
					return
				}
				logger.Debug("sent: %s", msg)
			case <-done:
				return
			}
		}
	}()

	// 协程2：接收客户端消息
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			logger.Debug("read error: %v", err)
			break
		}
		logger.Debug("recv: %s - %v", message, mt)

		// 可选：处理客户端消息
		// 例如：回复、触发某些动作等
	}

	// 主协程退出时通知定时发送协程退出
	close(done)
}
