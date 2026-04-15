package goeth

import (
	"fmt"

	"gorm.io/gorm"
)

func Run(db *gorm.DB) {
	fmt.Println("open goeth success")
}
