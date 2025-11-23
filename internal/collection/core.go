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
	if value == nil {
		*s = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan StringSlice: %v", value)
	}

	// 尝试 JSON 数组解析
	if bytes[0] == '[' {
		return json.Unmarshal(bytes, s)
	}

	// 普通字符串 → 转成切片
	*s = StringSlice{string(bytes)}
	return nil
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
	Id            int64    `json:"id" form:"id" gorm:"column:id"`
	Date          JSONTime `json:"date" form:"date" gorm:"column:updated_at"`
	Url           string   `json:"url" form:"url" gorm:"column:url"`
	InjectionType string   `json:"injectionType" form:"injectionType" gorm:"column:tag"`

	Directories *StringSlice `json:"directories" form:"directories" gorm:"column:directories"`
	Injection   *StringSlice `json:"injection" form:"injection" gorm:"column:injection"`
	Domains     *StringSlice `json:"domains" form:"domains" gorm:"column:domains"`
	Ports       *IntSlice    `json:"ports" form:"ports" gorm:"column:ports"`

	ManagerUrl  string `json:"managerUrl" form:"managerUrl"`
	ManagerUser string `json:"managerUser" form:"managerUser"`
	ManagerPass string `json:"managerPass" form:"managerPass"`
}

func InsertUrlInfo(
	url, injectionType, managerUrl, managerUser, managerPass string, injection StringSlice,
) error {
	// ... [校验 & 密码加密代码] ...

	// ✅ 关键：用 time.Now() 生成当前时间
	now := time.Now()

	record := &UrlInfo{
		Url:           url,
		InjectionType: injectionType,
		Injection:     &injection,
		ManagerUrl:    managerUrl,
		ManagerUser:   managerUser,
		ManagerPass:   managerPass,
		Date:          JSONTime(now),

		Directories: &StringSlice{},
		Domains:     &StringSlice{},
		Ports:       &IntSlice{},
	}

	ret := backend.DB.Create(record)
	if ret.Error != nil {
		logger.Error("insert failed: %v", ret.Error)
		return ret.Error
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
