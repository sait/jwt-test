package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var GlobalDB *gorm.DB

func InitDatabase() (err error) {
	GlobalDB, err = gorm.Open(sqlite.Open("auth.db"), &gorm.Config{})
	if err != nil {
		return
	}

	return
}
