package db

import (
	"fmt"
	"time"

	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"moul.io/zapgorm2"
)

type Database struct {
	Type            string
	Addr            string
	Port            string
	DBName          string
	Username        string
	Password        string
	MaxIdleConns    int           // 设置空闲连接池中连接的最大数量
	MaxOpenConns    int           // 设置打开数据库连接的最大数量
	ConnMaxLifetime time.Duration // 设置连接可复用的最大时间。
}

var logSugar *zap.SugaredLogger

// 设置gorm日志使用zap
var GormLogger zapgorm2.Logger

func (db *Database) GetGormDB(zaplog *zap.Logger) (*gorm.DB, error) {
	return db.NewDBConnect(zaplog)
}

// NewDBConnect 获取gorm.DB
func (db *Database) NewDBConnect(zaplog *zap.Logger) (GormDB *gorm.DB, err error) {
	// 使用zap 接收gorm日志
	GormLogger = zapgorm2.New(zaplog)
	GormLogger.SetAsDefault()
	logSugar = zaplog.Sugar()

	var dialector gorm.Dialector
	switch db.Type {
	case "sqlite":
		dialector = sqlite.Open(db.DBName)

	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			db.Username,
			db.Password,
			db.Addr,
			db.Port,
			db.DBName,
		)
		dialector = mysql.New(mysql.Config{
			DSN:                       dsn,
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		})
	case "postgresql":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			db.Addr,
			db.Username,
			db.Password,
			db.DBName,
			db.Port,
		)
		dialector = postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage,
		})
	default:
		logSugar.Error("The database is not supported, please choice [sqlite],[mysql] or [postgresql]")
		return
	}
	logSugar.Debugf("使用的数据库类型是: %s", db.Type)
	GormDB, err = gorm.Open(dialector, &gorm.Config{
		Logger:                                   GormLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 设置创建表名时不使用复数
		},
	})
	if db.MaxOpenConns != 0 {
		// Gorm 使用database/sql 维护连接池
		sqlDB, _ := GormDB.DB()
		// 设置空闲连接池中连接的最大数量
		sqlDB.SetMaxIdleConns(db.MaxIdleConns)
		// 设置打开数据库连接的最大数量
		sqlDB.SetMaxOpenConns(db.MaxOpenConns)
		// 设置了连接可复用的最大时间
		sqlDB.SetConnMaxLifetime(db.ConnMaxLifetime)
	} else {
		sqlDB, _ := GormDB.DB()
		// SetMaxIdleConns 设置空闲连接池中连接的最大数量
		sqlDB.SetMaxIdleConns(10)
		// SetMaxOpenConns 设置打开数据库连接的最大数量。
		sqlDB.SetMaxOpenConns(100)
		// SetConnMaxLifetime 设置连接可复用的最大时间。
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	return
}
