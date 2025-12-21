package main

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	//"gorm.io/driver/mysql"
	//"gorm.io/gorm"
)

func main() {
	var dsn string = "root:a123456@tcp(127.0.0.1:3306)/renren_fast?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//Run_sql1(db)
	//Run_sql2(db)
	Run_advanced(db)

	//db, err := sqlx.Open("mysql", dsn)
	//if err != nil {
	//	panic(err)
	//}
	//Run_sqlx(db)
}
