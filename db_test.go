package db

import (
	"testing"

	"github.com/serialt/sugar"
)

func TestGetGormDB(t *testing.T) {
	mydb := &Database{
		Type:     "mysql",
		Addr:     "10.0.16.10",
		Port:     "3306",
		DBName:   "exmail",
		Username: "root",
		Password: "rocky",
	}
	_, err := mydb.NewDBConnect(sugar.NewLogger("debug", "", "", false))
	// DB.AutoMigrate(&Department{})
	// DB.AutoMigrate(&Userlist{})
	if err != nil {
		t.Errorf("DB connect failed: %v", err)
	} else {
		t.Logf("DB connect succeed: %v", err)
	}
}
