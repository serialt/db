# 连接数据库

### 下载库
```
go get github.com/serialt/db
```

### 使用方法
```
package main

import (
	"fmt"

	"github.com/serialt/db"
	"github.com/serialt/sugar"

	// "github.com/serialt/sugar"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	mydb := &db.Database{
		Type:     "mysql",
		Addr:     "10.0.16.10.",
		Port:     "3306",
		DBName:   "exmail",
		Username: "root",
		Password: "rocky",
	}
	sugar.SetLog("info", "")
	DB = mydb.NewDBConnect(sugar.Log)
	// DB.AutoMigrate(&Department{})
	// DB.AutoMigrate(&Userlist{})
	sugar.Info("db connect")
}
func main() {
	fmt.Print("msg")
	sugar.Info("db connect")
}

```

