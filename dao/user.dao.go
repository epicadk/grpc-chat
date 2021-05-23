package dao

import (
	"github.com/epicadk/grpc-chat/db"
	"github.com/epicadk/grpc-chat/mappers"
	"github.com/epicadk/grpc-chat/models"
)

type UserDao struct{}

func (ud *UserDao) Create(user *models.User) error {
	u := mappers.UserProtoToDB(user)
	err := u.SaveToDB()

	return err
}

func (ud *UserDao) FindByUsername(phonenumber string) (*db.User, error) {
	user := db.User{
		Phonenumber: phonenumber,
	}

	err := user.FindUser()
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ud *UserDao) CheckCredentials(phonenumber, password string) (*db.User, error) {
	user := db.User{
		Phonenumber: phonenumber,
	}
	if err := user.CheckPassword(password); err != nil {
		return nil, err
	}
	return &user, nil
}
