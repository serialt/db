package db

import (
	"testing"

	"github.com/serialt/sugar/v3"
)

func TestGetGormDB(t *testing.T) {

	type User struct {
		Name string
		Age  int
	}
	gslog := sugar.New(sugar.WithLevel("debug"))
	gslog.Debug("hello")
	sqlDB, err := New("test.db", gslog)
	if err != nil {
		t.Fatal("open db failed")
	}
	sqlDB.AutoMigrate(&User{})

	jerry := &User{Name: "jerry", Age: 18}

	sqlDB.Model(&User{}).Debug().Create(&jerry)

	newJerry := &User{}

	sqlDB.Model(&User{}).Where("name = ?", "jerry").First(&newJerry)
	t.Logf("Get user from db, data=%v", newJerry)
}
