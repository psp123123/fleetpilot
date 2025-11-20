package backend

import (
	"fleetpilot/common/config"
	"fleetpilot/common/logger"
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm/schema"

	"gorm.io/gorm"
)

// 关于mysql的相关操作
var (
	DB   *gorm.DB
	once sync.Once
)

// 初始化mysql
func InitDB() error {

	once.Do(func() {
		// 拼接连接信息
		m := config.GlobalCfg.Mysql
		logger.Debug("------%v", m)
		dsn := fmt.Sprintf(
			"%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			m.Username,
			m.Password,
			m.Address,
			m.Dbname,
		)
		var err error
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 使用单数表名，启用此选项，此时 `UrlInfo` 的表名会是 `url_info`
			},
		})
		if err != nil {
			logger.Error("init db error:", err)

		}

		// 配置连接池
		sqlDB, err := DB.DB()
		if err != nil {
			logger.Error("get sql.DB error", err)

		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		logger.Info("mysql init success")
	})
	return nil
}

// 从mysql读取单条记录，传入任意参数，返回单条map
func GetMysqlOneData(queryTable string, conDiction map[string]interface{}) (*User, error) {

	//data := make(map[string]interface{})

	var user User

	result := DB.Table(queryTable).Where(conDiction).Take(&user)
	if result.Error != nil {
		logger.Error("查询失败:", result.Error)
		return nil, result.Error
	}

	logger.Debug("get user info ,userID column is %v", user.UserID)
	return &user, nil
}

type User struct {
	ID           int64     `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	UserID       string    `gorm:"column:userID;type:char(32);not null;comment:用户唯一ID" json:"userID"`
	Username     string    `gorm:"type:varchar(50);not null;unique;comment:用户名" json:"username"`
	PasswordHash string    `gorm:"type:varchar(255);not null;comment:密码哈希" json:"password_hash"`
	Email        string    `gorm:"type:varchar(100);not null;unique;comment:邮箱" json:"email"`
	Phone        string    `gorm:"type:varchar(20);index;comment:手机号" json:"phone"`
	Nickname     string    `gorm:"type:varchar(50);comment:昵称" json:"nickname"`
	Status       int8      `gorm:"type:tinyint(4);default:1;index;comment:状态(1=启用,0=禁用)" json:"status"`
	Role         string    `gorm:"type:varchar(30);default:'user';comment:角色" json:"role"`
	CreatedAt    time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	Avatar       string    `gorm:"type:varchar(255)" json:"avatar"`
	UpdatedAt    time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
}

func (User) TableName() string {
	return "user"
}
