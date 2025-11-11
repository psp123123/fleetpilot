package api

// websocket业务工具分发程序
import (
	"fleetpilot/common/logger"
	utils "fleetpilot/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// 定义工具处理客户端传入URL处理过程的接口
type ToolHandler interface {
	// 获取工具名称，用户注册和路由
	GetToolName() string
	Executed(conn *websocket.Conn, msg []byte) (interface{}, error)
}

// 解析http协议携带的参数
type ClientHttpStruct struct {
	ToolName string `form:"tool"`
	UserName string `form:"user"`
	Token    string `form:"token"`
}

// 解析函数：从gin上下文的查询参数中绑定数据到结构体
func (c *ClientHttpStruct) GetToolName() string {
	return c.ToolName
}

// 接收者为*ClientHttpStruct，返回绑定是否成功
func (c *ClientHttpStruct) ClientBindQuery(ctx *gin.Context) error {
	return ctx.ShouldBindQuery(c)
}

func WsHandler(ctx *gin.Context) {
	// 协议升级前获取参数
	var chs ClientHttpStruct

	// 调用接收者方法绑定查询参数
	if err := chs.ClientBindQuery(ctx); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"timestamp": time.Now().UnixMilli(),
			"msg":       "bad request" + err.Error(),
			"extra1":    "",
			"extra2":    gin.H{},
		})
		return
	}
	logger.Debug("升级协议前获取参数toolname:%v", chs.ToolName)
	// 验证用户信息
	// user := cw.UserName
	accessToken := chs.Token
	claims, err := utils.VerifyAccessToken(accessToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"timestamp": time.Now().UnixMilli(),
			"msg":       "Bad Request",
			"extra1":    "",
			"extra2":    gin.H{},
		})
	}
	logger.Debug("get token's claims is %v", claims)

	// 获取工具处理器GetHandlerManager
	manager := GetHandlerManager()
	handler, exsit := manager.GetHandler(chs.ToolName)
	if !exsit {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"timestamp": time.Now().UnixMilli(),
			"msg":       "Bad tool",
			"extra1":    "",
			"extra2":    gin.H{},
		})
		return
	}
	logger.Debug("get handler is %v", handler.GetToolName())
	// 升级协议
	c, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Error("ws协议升级失败", err)
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{
			"timestamp": time.Now().UnixMilli(),
			"msg":       err.Error(),
			"extra1":    "",
			"extra2":    gin.H{},
		})
		return
	}
	defer c.Close()

	// 接受客户端消息参数，并启用协程处理，并时刻返回
	for {
		_, msgBytes, err := c.ReadMessage()
		if err != nil {
			logger.Error("read client msg error: %v", err)
			break
		}

		logger.Debug("get client ws message is %s", string(msgBytes))
	}

}
