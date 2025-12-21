package main

import (
	"gorm.io/gorm"
)

// SQL语句练习
// 题目1：基本CRUD操作
// 假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
// 要求 ：
// 编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
// 编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
// 编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
// 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
func Run_sql1(db *gorm.DB) {
	//db.Debug().AutoMigrate(&Student{})
	// CREATE TABLE `students` (`id` bigint AUTO_INCREMENT,`name` varchar(20),`age` tinyint,`grade` varchar(20),PRIMARY KEY (`id`))
	//
	//var s1 Student = Student{Name: "张三", Age: 20, Grade: "三年级"}
	//var s2 Student = Student{Name: "李四", Age: 13, Grade: "一年级"}
	//db.Debug().Create(&s1)
	//db.Debug().Create(&s2)
	// INSERT INTO `students` (`name`,`age`,`grade`) VALUES ('张三',20,'三年级')
	// INSERT INTO `students` (`name`,`age`,`grade`) VALUES ('李四',13,'一年级')
	//
	//var sSlice []Student
	//db.Debug().Raw("SELECT * FROM students WHERE age > ?", 18).Scan(&sSlice)
	//fmt.Println(sSlice)
	//
	//result1 := db.Debug().Exec("UPDATE students SET grade = ? WHERE name = ?", "四年级", "张三")
	//if result1.Error != nil {
	//	panic(result1.Error)
	//}
	//if result1.RowsAffected == 0 {
	//	fmt.Println("没有满足条件的数据")
	//}
	//
	//result2 := db.Debug().Exec("DELETE FROM students WHERE age < ?", 15)
	//if result2.Error != nil {
	//	panic(result2.Error)
	//}
	//if result2.RowsAffected == 0 {
	//	fmt.Println("没有满足条件的数据")
	//}
}

type Student struct {
	ID    int64
	Name  string `gorm:"type:varchar(20)"`
	Age   int8
	Grade string `gorm:"type:varchar(20)"`
}
