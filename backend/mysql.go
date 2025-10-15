package backend

import (
	"fleetpilot/common/config"
	"fleetpilot/common/logger"
	"fmt"
	"time"

	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

// 关于mysql的相关操作

// 初始化mysql
func InitDB() (*gorm.DB, error) {
	m := config.GlobalCfg.Mysql.Mysqler

	// 拼接连接信息
	dsn := fmt.Sprintf(
		"%s:%s@(%s/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.Username,
		m.Password,
		m.Address,
		m.Dbname,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("init db error:", err)
		return nil, err
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("get sql.DB error", err)
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Info("mysql init success")
	return db, nil
}

// 从mysql读取单条记录，传入任意参数，返回单条map
func GetMysqlOneData(queryTable string, conDiction map[string]interface{}) map[string]interface{} {

	//data := make(map[string]interface{})
	db, err := InitDB()
	if err != nil {
		logger.Error("mysql初始化失败:", err)
	}

	data := make(map[string]interface{})

	result := db.Table(queryTable).Where(conDiction).Take(&data)
	if result.Error != nil {
		logger.Error("查询失败:", result.Error)
		return nil
	}

	return data
}
