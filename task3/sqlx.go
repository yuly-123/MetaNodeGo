package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

// Sqlx入门
// 题目1：使用SQL扩展库进行查询
// 假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
// 要求 ：
// 编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
// 编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
func Run_sqlx(db *sqlx.DB) {
	var eSlice []Employee
	err1 := db.Select(&eSlice, "SELECT * FROM employees WHERE department = ?", "技术部")
	if err1 != nil {
		panic(err1)
	}
	fmt.Println(eSlice)

	var e Employee
	err2 := db.Get(&e, "SELECT * FROM employees ORDER BY salary DESC LIMIT 1")
	if err2 != nil {
		panic(err2)
	}
	fmt.Println(e)
}

type Employee struct {
	ID         int64
	Name       string
	Department string
	Salary     int64
}
