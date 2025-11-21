package collection

import (
	"fleetpilot/backend"
	"fleetpilot/common/logger"
	"time"
)

type UrlInfo struct {
	Id            int64     `json:"id"  form:"form" gorm:"column:id"`
	Date          time.Time `json:"date" form:"date" gorm:"column:updated_at"`
	Url           string    `json:"url" form:"url" gorm:"column:url"`
	InjectionType string    `json:"injectionType"  form:"injectionType" gorm:"column:tag"`
	InjectionPath []string  `json:"injectionPath"  form:"injectionPath" gorm:"column:directories"`
}

func (c *UrlInfo) InsertData(id int64, date time.Time, url, injectionType string, injectionPath []string) error {
	record := &UrlInfo{
		id, date, url, injectionType, injectionPath,
	}

	ret := backend.DB.Create(record)
	if ret.Error != nil {
		logger.Error("insert failed :%v", ret.Error)
	}
	return nil
}

// get data from mysql
func (c *UrlInfo) GetData(url string) error {

	ret := backend.DB.Where("url = ?", url).First(c)
	if ret.Error != nil {
		logger.Error("get %s from mysql failed :%v", url, ret.Error)
	}
	return nil
}

// get all data from mysql
func (c *UrlInfo) GetAllData() ([]UrlInfo, error) {
	var urlListData []UrlInfo
	ret := backend.DB.Find(&urlListData)

	if ret.Error != nil {
		logger.Error("获取url列表失败:%v", ret.Error)
		return nil, ret.Error
	}
	logger.Debug("get all data :%v", urlListData)
	return urlListData, nil
}

// // 返回结构体
// func GetCollectionStruct() *CollectionInfo {
// 	return &CollectionInfo{}
// }
