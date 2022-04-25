package gorm_cache

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	username     = ""
	password     = ""
	databaseName = ""
	ip           = ""
	port         = ""
)

func log(format string, a ...interface{}) {
	timeStr := time.Now().Format("2006-01-02 15:04:05.999")
	fmt.Printf(timeStr+" "+format+"\n", a...)
}

// here we only test with mysql
func main() {
	var err error
	defer func() {
		if err != nil {
			log("integration exits with error: %v", err)
		}
	}()

	log("integration test start...")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, ip, port, databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
	})
	if err != nil {
		log("open db error: %v", err)
		return
	}

	log("integration test ends.")
}
