package collection

import (
	"fleetpilot/common/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 路由注册
type CollectionHandler struct{}

func (h *CollectionHandler) RegisterRoutes(rg *gin.RouterGroup) {
	logger.Debug("----enter RegisterRoutes")
	collection := rg.Group("/collection")
	{
		collection.POST("/urlcreate", h.InfoCreate)
		collection.PATCH("/urlpatch", h.InfoPatch)
		collection.GET("/urlget", h.GetData)
	}
}

// 客户端信息写入mysql
func (h *CollectionHandler) InfoCreate(c *gin.Context) {
	// bind struct
	var Info UrlInfo
	if err := c.ShouldBind(&Info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// insert data to mysql
	if err := Info.InsertData(Info.Url, Info.InjectionType, Info.InjectionPath); err != nil {
		logger.Error("data writed failed :%V", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	// response newest data to client
	if err := Info.GetData(Info.Url); err != nil {
		logger.Error("get data failed :%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": Info})
	}
}

func (h *CollectionHandler) InfoPatch(c *gin.Context) {}

func (h *CollectionHandler) GetData(c *gin.Context) {
	var Info UrlInfo
	// response newest data to client
	if err := Info.GetData(Info.Url); err != nil {
		logger.Error("get data failed :%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": Info})
	}
}

func NewCollectionHandler() *CollectionHandler {
	return &CollectionHandler{}
}
