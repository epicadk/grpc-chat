package db

import "errors"

type User struct {
	Id       string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username string `gorm:"uniqueindex;unique;not null"`
	Password string `gorm:"not null"`
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
	if user.Password != password {
		return errors.New("passwords do not match")
	}
	return nil

}
