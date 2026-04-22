package main

import (
	// "MetaNodeGo/lesson04"
	// "MetaNodeGo/goeth1"
	// "MetaNodeGo/goeth2"
	// "MetaNodeGo/goeth3"
	// "MetaNodeGo/goeth4"
	// "MetaNodeGo/goeth5"
	// "MetaNodeGo/goeth6"
	"MetaNodeGo/goeth7"

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
	// goeth1.Run(db)
	// goeth2.Run(db)
	// goeth3.Run(db)
	// goeth4.Run(db)
	// goeth5.Run(db)
	// goeth6.Run(db)
	goeth7.Run(db)
}
