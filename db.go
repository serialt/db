package db

import (
	"github.com/glebarez/sqlite"

	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func New(file string) *gorm.DB {

	gromdb, err := gorm.Open(sqlite.Open(file), &gorm.Config{

		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 设置创建表名时不使用复数
		},
	})
	if err != nil {
		slog.Info("open db failed", "error", err)
	}

	return gromdb

}
