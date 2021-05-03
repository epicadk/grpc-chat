package db

import (
	"github.com/epicadk/grpc-chat/utils"
	"gorm.io/gorm"
)

type User struct {
	// Primary key is indexed by default
	Phonenumber string `gorm:"primaryKey"`
	DisplayName string `gorm:"not null"`
	Password    string `gorm:"not null"`
}

// GORM hook, validate data here.
func (user *User) BeforeCreate(tx *gorm.DB) error {
	hashedpass, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
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
	// TODO user hash
	if err != nil {
		return err
	}
	return utils.ComparePassword(user.Password, password)
}
