package dao

import (
	"errors"

	"github.com/epicadk/grpc-chat/db"
	"github.com/epicadk/grpc-chat/models"
	"github.com/epicadk/grpc-chat/utils"
	"gorm.io/gorm"
)

type UserDao struct{}

func (ud *UserDao) SaveUser(user *models.User) (string, error) {
	u := utils.UserProtoToDb(user)
	tx := db.DBconn.Create(u)
	if tx.Error != nil {
		// TODO throw custom errors
		return "", tx.Error
	}
	return u.Id, nil
}

func (ud *UserDao) FindUserByUsername(username string) (*db.User, error) {
	var user db.User
	tx := db.DBconn.Where("username :=", username).First(&user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			// TODO throw custom error
			return nil, errors.New("no user found")
		}
		return nil, tx.Error
	}
	return &user, nil
}

func (ud *UserDao) FindUserByID(userid string) (*db.User, error) {
	var user db.User
	tx := db.DBconn.Where("Id :=", userid).First(&user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			// TODO throw custom error
			return nil, errors.New("no user found")
		}
		return nil, tx.Error
	}
	return &user, nil
}
