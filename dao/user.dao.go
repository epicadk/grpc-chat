package dao

import (
	"github.com/epicadk/grpc-chat/db"
	"github.com/epicadk/grpc-chat/models"
	"github.com/epicadk/grpc-chat/utils"
)

type UserDao struct{}

// TODO some validation
func (ud *UserDao) Create(user *models.User) error {
	u := utils.UserProtoToDb(user)
	err := u.SaveToDB()
	user.UserID = u.Id

	return err
}

func (ud *UserDao) FindByUsername(username string) (*db.User, error) {
	user := db.User{
		Username: username,
	}

	err := user.FindUser()
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ud *UserDao) CheckCredentials(username, password string) error {
	user := db.User{
		Username: username,
	}
	return user.CheckPassword(password)

}
