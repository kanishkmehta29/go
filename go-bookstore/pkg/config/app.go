package config

import(
   "fmt"
   "gorm.io/driver/mysql"
   "gorm.io/gorm"
)

var db *gorm.DB

func Connect() {
    dsn := "root:root_password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
    var err error
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    fmt.Println("Successfully connected to MySQL database")
}

func GetDB() *gorm.DB{
   return db
}