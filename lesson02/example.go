package lesson02

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name       string
	Age        int
	Birthday   time.Time
	CreditCard CreditCard
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Name = u.Name + "_123123"
	return
}

type CreditCard struct {
	gorm.Model
	Number string
	UserID uint
}

//type BankUser struct {
//	gorm.Model
//	Name       string
//	CreditCard CreditCard
//}

func Run(db *gorm.DB) {
	//db.AutoMigrate(&User{})
	//db.AutoMigrate(&BankUser{})
	//db.AutoMigrate(&CreditCard{})

	//users := []User{
	//	{CreditCard: CreditCard{Number: "003"}, Name: "张三", Age: 31, Birthday: time.Now()},
	//	{CreditCard: CreditCard{Number: "004"}, Name: "李四", Age: 32, Birthday: time.Now()},
	//	{CreditCard: CreditCard{Number: "005"}, Name: "王五", Age: 33, Birthday: time.Now()},
	//}
	//result := db.Session(&gorm.Session{SkipHooks: true}).Create(users)
	//fmt.Println(result.RowsAffected)

	var users []User
	db.Raw("SELECT * FROM users").Scan(&users)
	fmt.Println(len(users))

	var maps []map[string]interface{}
	db.Raw("SELECT * FROM users").Scan(&maps)

	for i, row := range maps {
		fmt.Printf("第 %d 行:\n", i+1)
		for key, value := range row {
			fmt.Printf("  %s: %s\n", key, value)
		}
		fmt.Println()
	}

	//var user User
	//user.ID = 10
	//result := db.Debug().First(&user, "name = ?", "张三")
	//db.Debug().Take(&user)
	//db.Debug().Last(&user)
	//fmt.Println(result)
	//if errors.Is(result.Error, gorm.ErrRecordNotFound) {
	//	fmt.Println("=============>", result.Error)
	//}
}
