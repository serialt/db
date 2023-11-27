package db

import (
	"log/slog"
	"testing"

	"github.com/serialt/sugar/v3"
)

func TestGetGormDB(t *testing.T) {
	slog.SetDefault(sugar.New())
	mydb := &Database{
		Type:     "mysql",
		Addr:     "10.0.16.10",
		Port:     "3336",
		DBName:   "mysql",
		Username: "root",
		Password: "rocky",
	}

	_, err := mydb.NewDBConnect()
	// DB.AutoMigrate(&Department{})
	// DB.AutoMigrate(&Userlist{})
	if err != nil {
		t.Errorf("DB connect failed: %v", err)
	} else {
		t.Logf("DB connect succeed: %v", err)
	}
}
