package services

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string
	Lastname     string
	Email        string `gorm:"unique"`
	Phone        string `gorm:"unique"`
	IDCardNumber string `gorm:"unique"`
	Role         string
	Password     string
}

type Role struct {
	ID    uint `gorm:"primarykey"`
	Title string
}

func ApplyMigrations(db *gorm.DB) {
	err := db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}
}
