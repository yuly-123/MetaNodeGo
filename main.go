package main

import (
	// "MetaNodeGo/lesson04"
	// "MetaNodeGo/goeth1"
	"MetaNodeGo/goeth2"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	var dsn string = "root:a123456@tcp(127.0.0.1:3306)/renren_fast?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// lesson04.Run(db)
	goeth2.Run(db)
}
