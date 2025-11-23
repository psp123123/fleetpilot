package collection

import (
	"fleetpilot/common/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 路由注册
type CollectionHandler struct{}

func (h *CollectionHandler) RegisterRoutes(rg *gin.RouterGroup) {
	logger.Debug("----enter RegisterRoutes")
	collection := rg.Group("/collection")
	{
		collection.DELETE("/:id", h.InfoDelete)
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
	if err := InsertUrlInfo(Info.Url, Info.InjectionType, Info.Injection, Info.ManagerUrl, Info.ManagerUser, Info.ManagerPass); err != nil {
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

func (h *CollectionHandler) InfoDelete(c *gin.Context) {
	var Info UrlInfo
	// 获取删除的数据id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("delete id is invalid: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
	}

	err = Info.DeleteData(id)
	if err != nil {
		logger.Error("invalid id", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "delete ok",
	})
}
func (h *CollectionHandler) GetData(c *gin.Context) {
	var Info UrlInfo
	// response newest data to client
	urllist, err := Info.GetAllData()
	if err != nil {
		logger.Error("get data failed :%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"urllist": urllist})
	}
}

func NewCollectionHandler() *CollectionHandler {
	return &CollectionHandler{}
}
