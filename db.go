package db

import (
	"log/slog"

	"github.com/glebarez/sqlite"

	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func New(file string, gslog *slog.Logger) (gromdb *gorm.DB, err error) {
	gormLogger := slogGorm.New(
		slogGorm.WithHandler(gslog.Handler()),
		slogGorm.WithTraceAll(), // trace all messages
		slogGorm.SetLogLevel(slogGorm.DefaultLogType, slog.LevelDebug),
	)

	gromdb, err = gorm.Open(sqlite.Open(file), &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 设置创建表名时不使用复数
		},
	})
	return

}
