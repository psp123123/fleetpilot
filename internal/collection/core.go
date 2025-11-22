package collection

import (
	"database/sql/driver"
	"encoding/json"
	"fleetpilot/backend"
	"fleetpilot/common/logger"
	"fmt"
	"time"
)

// mysql中数据字符串转换为整数切片类型
type IntSlice []int

// 实现sql.Scanner接口
func (i *IntSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		// 如果value不是[]byte类型（例如NULL），尝试处理
		if value == nil {
			*i = nil
			return nil
		}
		return fmt.Errorf("failed to scan IntSlice: %v, type: %T", value, value)
	}

	// 直接将JSON数组解析到IntSlice（[]int）中
	return json.Unmarshal(bytes, i)
}

// 实现driver.Valuer接口
func (i *IntSlice) Value() (driver.Value, error) {
	// 如果切片为nil或空，返回nil，数据库中会存储为NULL
	if i == nil || len(*i) == 0 {
		return nil, nil
	}
	return json.Marshal(i)
}

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

func (t JSONTime) Value() (driver.Value, error) {
	ts := time.Time(t)
	return ts.Format("2006-01-02 15:04:05"), nil
}

func (t *JSONTime) Scan(value interface{}) error {
	if value == nil {
		*t = JSONTime(time.Time{})
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*t = JSONTime(v)
	case []byte:
		parsed, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}
		*t = JSONTime(parsed)
	case string:
		parsed, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
		*t = JSONTime(parsed)
	default:
		return fmt.Errorf("unsupported type %T", v)
	}

	return nil
}

// 定义数据查询
type UrlInfo struct {
	Id            int64       `json:"id"  form:"form" gorm:"column:id"`
	Date          JSONTime    `json:"date" form:"date" gorm:"column:updated_at"`
	Url           string      `json:"url" form:"url" gorm:"column:url"`
	InjectionType string      `json:"injectionType"  form:"injectionType" gorm:"column:tag"`
	InjectionPath StringSlice `json:"injectionPath"  form:"injectionPath" gorm:"column:directories"`
	Injection     string      `json:"injection"  form:"injection" gorm:"column:injection"`
	Domains       StringSlice `json:"domains" form:"domains" gorm:"column:domains"`
	Ports         IntSlice    `json:"ports" form:"ports" gorm:"column:ports;serializer:json"`
	ManagerUrl    string      `json:"managerUrl" form:"managerUrl"`   // gorm:"column:manager_url;serializer:json"`
	ManagerUser   string      `json:"managerUser" form:"managerUser"` //gorm:"column:manager_user;serializer:json"`
	ManagerPass   string      `json:"managerPass" form:"managerPass"` //gorm:"column:manager_pass;seiralizer:json"`
}

func (c *UrlInfo) InsertData(url, injectionType, injection, ManagerUrl, ManagerUser, ManagerPass string) error {
	//record := &UrlInfo{}

	ret := backend.DB.Create(c)
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

// delete one data
func (c *UrlInfo) DeleteData(id int64) error {
	err := backend.DB.Delete(c, id).Error
	if err != nil {
		logger.Error("delete data failed :%v", err)
		return err
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
