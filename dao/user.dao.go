package dao

import (
	"github.com/epicadk/grpc-chat/db"
	"github.com/epicadk/grpc-chat/models"
	"github.com/epicadk/grpc-chat/utils"
)

type UserDao struct{}

// TODO some validation
func (ud *UserDao) Create(user *models.User) error {
	u := utils.UserProtoToDB(user)
	err := u.SaveToDB()
	user.UserID = u.ID

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

func (ud *UserDao) CheckCredentials(phonenumber, password string) error {
	user := db.User{
		Phonenumber: phonenumber,
	}
	return user.CheckPassword(password)

}
