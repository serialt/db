package db

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/glebarez/sqlite"
	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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

// 设置gorm日志使用zap

// NewDBConnect 获取gorm.DB
func NewDBConnect(conf *Database, gslog *slog.Logger) (GormDB *gorm.DB, err error) {
	gormLogger := slogGorm.New(
		slogGorm.WithHandler(gslog.Handler()),
		slogGorm.WithTraceAll(), // trace all messages
		slogGorm.SetLogLevel(slogGorm.DefaultLogType, slog.LevelDebug),
	)

	var dialector gorm.Dialector
	switch conf.Type {
	case "sqlite":
		dialector = sqlite.Open(conf.DBName)

	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			conf.Username,
			conf.Password,
			conf.Addr,
			conf.Port,
			conf.DBName,
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
			conf.Addr,
			conf.Username,
			conf.Password,
			conf.DBName,
			conf.Port,
		)
		dialector = postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage,
		})
	default:
		return nil, errors.New("The database is not supported, please choice [sqlite],[mysql] or [postgresql]")
	}
	GormDB, err = gorm.Open(dialector, &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 设置创建表名时不使用复数
		},
	})
	if conf.MaxOpenConns != 0 {
		// Gorm 使用database/sql 维护连接池
		sqlDB, _ := GormDB.DB()
		// 设置空闲连接池中连接的最大数量
		sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
		// 设置打开数据库连接的最大数量
		sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
		// 设置了连接可复用的最大时间
		sqlDB.SetConnMaxLifetime(conf.ConnMaxLifetime)
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
