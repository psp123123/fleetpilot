package collection

import (
	"database/sql/driver"
	"encoding/json"
	"fleetpilot/backend"
	"fleetpilot/common/logger"
	"fmt"
	"time"
)

// mysql中数据字符串转换切片类型
type StringSlice []string

// 实现sql.Scanner接口
func (s *StringSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("falied to scan StringSlice: %v", value)
	}

	return json.Unmarshal(bytes, s)
}

// 实现driver.valuer接口
func (s *StringSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// 自定义时间格式
type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

// 定义数据查询
type UrlInfo struct {
	Id            int64       `json:"id"  form:"form" gorm:"column:id"`
	Date          JSONTime    `json:"date" form:"date" gorm:"column:updated_at"`
	Url           string      `json:"url" form:"url" gorm:"column:url"`
	InjectionType string      `json:"injectionType"  form:"injectionType" gorm:"column:tag"`
	InjectionPath StringSlice `json:"injectionPath"  form:"injectionPath" gorm:"column:directories"`
}

func (c *UrlInfo) InsertData(id int64, date JSONTime, url, injectionType string, injectionPath []string) error {
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
