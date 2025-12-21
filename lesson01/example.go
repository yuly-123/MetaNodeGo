package lesson01

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uint           // Standard field for the primary key
	Name         string         // A regular string field
	Email        *string        // A pointer to a string, allowing for null values
	Age          uint8          // An unsigned 8-bit integer
	Birthday     *time.Time     // A pointer to time.Time, can be null
	MemberNumber sql.NullString // Uses sql.NullString to handle nullable strings
	ActivatedAt  sql.NullTime   // Uses sql.NullTime for nullable time fields
	CreatedAt    time.Time      // Automatically managed by GORM for creation time
	UpdatedAt    time.Time      // Automatically managed by GORM for update time
	ignored      string         // fields that aren't exported are ignored
}
type Member struct {
	gorm.Model
	Name string
	Age  uint8
}
type Author struct {
	Name  string
	Email string
}

type Blog struct {
	Author
	ID      int
	Upvotes int32
}

type Blog2 struct {
	ID      int64
	Author  Author `gorm:"embedded;embeddedPrefix:author_"`
	Upvotes int32  `gorm:"column:votes"`
}

func Run(db *gorm.DB) {
	db.AutoMigrate(&User{})
	//db.AutoMigrate(&Member{})
	//db.AutoMigrate(&Blog{})
	db.AutoMigrate(&Blog2{})

	user := &User{}
	user.MemberNumber.Valid = true
	result := db.Create(user)
	fmt.Println(result.RowsAffected)
}
