package lesson04

import (
	"fmt"
	"gorm.io/gorm"
)

type Dog struct {
	ID   int
	Name string
	Toy  Toy `gorm:"polymorphic:Owner"`
}
type Cat struct {
	ID   int
	Name string
	Toy  Toy `gorm:"polymorphic:Owner"`
}
type Toy struct {
	ID        int
	Name      string
	OwnerType string
	OwnerID   int
}

func Run(db *gorm.DB) {
	//db.AutoMigrate(&Dog{}, &Cat{}, &Toy{})
	//
	//db.Create(&Dog{Name: "WangCai", Toy: Toy{Name: "gutou"}})
	//db.Create(&Cat{Name: "MiMi", Toy: Toy{Name: "doumaobang"}})

	var dog Dog
	//var cat Cat
	//db.Preload("Toy").First(&dog)
	//db.Preload("Toy").First(&cat)
	stmt := db.Session(&gorm.Session{DryRun: true}).First(&dog).Statement
	fmt.Println(stmt.SQL.String())
	fmt.Println(stmt.Vars)
	//fmt.Println(dog)
	//fmt.Println(cat)
}
