package main

import (
	"fmt"
	"gorm.io/gorm"
)

// SQL语句练习
// 题目2：事务语句
// 假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
// 要求 ：
// 编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。
// 在事务中，需要先检查账户 A 的余额是否足够，
// 如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。
// 如果余额不足，则回滚事务。
func Run_sql2(db *gorm.DB) {
	//db.Debug().AutoMigrate(&Account{}, &Transaction{})
	// CREATE TABLE `accounts` (`id` bigint AUTO_INCREMENT,`balance` bigint,PRIMARY KEY (`id`))
	// CREATE TABLE `transactions` (`id` bigint AUTO_INCREMENT,`from_account_id` bigint,`to_account_id` bigint,`amount` bigint,PRIMARY KEY (`id`))

	//var a Account = Account{Balance: 100}
	//var b Account = Account{Balance: 100}
	//db.Debug().Create(&a)
	//db.Debug().Create(&b)
	// INSERT INTO `accounts` (`balance`) VALUES (100)
	// INSERT INTO `accounts` (`balance`) VALUES (100)

	var a1 Account
	var b1 Account
	a1.ID = 1
	b1.ID = 2
	db.Debug().First(&a1)
	db.Debug().First(&b1)
	fmt.Println(a1)
	fmt.Println(b1)
	// SELECT * FROM `accounts` WHERE `accounts`.`id` = 1 ORDER BY `accounts`.`id` LIMIT 1
	// SELECT * FROM `accounts` WHERE `accounts`.`id` = 2 ORDER BY `accounts`.`id` LIMIT 1

	tx := db.Begin()

	if a1.Balance < 100 {
		fmt.Println("余额不足，转账失败！")
		tx.Rollback()
	}

	result1 := tx.Debug().Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", 100, a1.ID)
	if result1.Error != nil {
		panic(result1.Error)
		tx.Rollback()
	}
	result2 := tx.Debug().Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", 100, b1.ID)
	if result2.Error != nil {
		panic(result2.Error)
		tx.Rollback()
	}

	var t Transaction = Transaction{FromAccountId: a1.ID, ToAccountId: b1.ID, Amount: 100}
	tx.Debug().Create(&t)
	// INSERT INTO `transactions` (`from_account_id`,`to_account_id`,`amount`) VALUES (1,2,100)

	tx.Commit()
}

type Account struct {
	ID      int64
	Balance int64
}
type Transaction struct {
	ID            int64
	FromAccountId int64
	ToAccountId   int64
	Amount        int64
}
