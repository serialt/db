package db

import (
	"github.com/glebarez/sqlite"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func New(file string) (gromdb *gorm.DB, err error) {
	GormLoger := SetSlog()
	gromdb, err = gorm.Open(sqlite.Open(file), &gorm.Config{
		Logger:                                   GormLoger,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 设置创建表名时不使用复数
		},
	})
	return

}
