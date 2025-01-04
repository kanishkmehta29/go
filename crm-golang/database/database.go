package database

import (
	"gorm.io/gorm"
	_ "gorm.io/driver/postgres"
)

var DB *gorm.DB