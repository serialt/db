package db

import (
	"log/slog"
	"testing"
)

func TestGetGormDB(t *testing.T) {
	type Person struct {
		Name string
		Age  int64
	}
	var dbList []*Database
	mydb := &Database{
		Type:     "mysql",
		Addr:     "10.0.16.10",
		Port:     "3306",
		DBName:   "dbtest",
		Username: "root",
		Password: "centos",
	}
	sqlitedb := &Database{
		Type:   "sqlite",
		DBName: "dbtest.db",
	}
	pgdb := &Database{
		Type:     "postgresql",
		Addr:     "10.0.16.10",
		Port:     "5432",
		DBName:   "dbtest",
		Username: "postgres",
		Password: "centos",
	}
	dbList = append(dbList, sqlitedb)
	dbList = append(dbList, mydb)
	dbList = append(dbList, pgdb)

	for i, v := range dbList {
		db, err := NewDBConnect(v, slog.Default())
		if err != nil {
			t.Errorf("DB connect failed: %v", err)
		} else {
			t.Logf("DB connect succeed, index: %v", i)
		}

		db.AutoMigrate(&Person{})
		result := db.Model(&Person{}).Create(&Person{Name: "Jerry", Age: 18})
		if result.Error != nil {
			t.Errorf("create user failed: %v", err)
		}
		t.Logf("create user succeed")

		db.Exec("drop table person")
		if result.Error != nil {
			t.Errorf("drop table of user failed: %v", err)
		}
		t.Logf("drop table succeed")
	}

}
