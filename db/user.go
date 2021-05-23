package db

import (
	"log"

	"github.com/epicadk/grpc-chat/utils"
	"gorm.io/gorm"
)

type User struct {
	Phonenumber string `gorm:"primaryKey"` // Primary key is indexed by default
	DisplayName string `gorm:"not null"`   // Display name of the user
	Password    string `gorm:"not null"`   // Hashed password of the user
}

// GORM hook, validate data here.
func (user *User) BeforeCreate(tx *gorm.DB) error {
	hashedpass, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Fatal(err)
	}
	user.Password = hashedpass
	return err
}

func (user *User) SaveToDB() error {
	return DBconn.Create(user).Error
}

func (user *User) FindUser() error {
	return DBconn.Where(user).First(user).Error
}

func (user *User) CheckPassword(password string) error {
	err := DBconn.Where(user).First(user).Error
	if err != nil {
		return err
	}
	return utils.ComparePassword(user.Password, password)
}
