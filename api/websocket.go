package api

import (
	"fleetpilot/common/logger"
	utils "fleetpilot/utils"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// HTTP查询参数
type ClientHttpStruct struct {
	ToolName string `form:"tool"`
	UserName string `form:"user"`
	Token    string `form:"token"`
}

func WsHandler(ctx *gin.Context) {
	var chs ClientHttpStruct
	if err := ctx.ShouldBindQuery(&chs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "bad request"})
		return
	}

	logger.Debug("升级协议前获取参数 toolname: %v", chs.ToolName)

	// token 验证
	if _, err := utils.VerifyAccessToken(chs.Token); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "invalid token"})
		return
	}

	// 获取 handler
	manager := GetHandlerManager()
	handler, exists := manager.GetHandler(chs.ToolName)
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "unknown tool"})
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Error("ws协议升级失败", err)
		return
	}
	defer conn.Close()

	// 循环接收客户端消息
	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			logger.Error("read message error: %v", err)
			break
		}

		logger.Debug("ws recv: %s", msgBytes)

		// 执行对应工具
		if err := handler.Executed(conn, msgBytes); err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("error: "+err.Error()))
			break
		}
	}
}
