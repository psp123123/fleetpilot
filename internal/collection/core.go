package collection

import (
	"fleetpilot/backend"
	"fleetpilot/common/logger"
)

type UrlInfo struct {
	Url           string `json:"url" db:"url" form:"url"`
	InjectionType string `json:"injectionType" db:"injectionType" form:"injectionType"`
	InjectionPath string `json:"injectionPath" db:"injectionPath" form:"injectionPath"`
}

func (c *UrlInfo) InsertData(url, injectionType, injectionPath string) error {
	record := &UrlInfo{
		url, injectionType, injectionPath,
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

// // 返回结构体
// func GetCollectionStruct() *CollectionInfo {
// 	return &CollectionInfo{}
// }
