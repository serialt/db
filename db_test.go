package db

import (
	"testing"
)

func TestGetGormDB(t *testing.T) {

	type User struct {
		Name string
		Age  int
	}

	sqlDB, err := New("test.db")
	if err != nil {
		t.Fatal("open db failed")
		panic(err)
	}
	sqlDB.AutoMigrate(&User{})

	jerry := &User{Name: "jerry", Age: 18}

	sqlDB.Model(&User{}).Create(&jerry)

	newJerry := &User{}

	sqlDB.Model(&User{}).Where("name = ?", "jerry").First(&newJerry)
	t.Logf("Get user from db, data=%v", newJerry)
}
