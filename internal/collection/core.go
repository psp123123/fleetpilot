package collection

import (
	"fleetpilot/backend"
	"fleetpilot/common/logger"
)

type UrlInfo struct {
	Date          string `json:"date" db:"updated_at" form:"date"`
	Url           string `json:"url" db:"url" form:"url"`
	InjectionType string `json:"injectionType" db:"tag" form:"injectionType"`
	InjectionPath string `json:"injectionPath" db:"directories" form:"injectionPath"`
}

func (c *UrlInfo) InsertData(date, url, injectionType, injectionPath string) error {
	record := &UrlInfo{
		date, url, injectionType, injectionPath,
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
