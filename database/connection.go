package database

import (
    "go-auth/models"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
    conn, err := gorm.Open(mysql.Open("zam0x1ma:rRd8)tHv_F]b?hVQ:s2lox@/go_auth"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    DB = conn

    conn.AutoMigrate(&models.User{})
}
