package dao

import (
	"errors"

	"github.com/epicadk/grpc-chat/db"
	"github.com/epicadk/grpc-chat/models"
	"github.com/epicadk/grpc-chat/utils"
	"gorm.io/gorm"
)

type ChatDao struct{}

func (cd ChatDao) SaveChat(chat *models.Message) error {

	tx := db.DBconn.Create(utils.ChatProtoToD(chat))
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (cd ChatDao) DeleteChat(chat *db.Chat) error {
	tx := db.DBconn.Delete(chat)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (cd ChatDao) FindChat(rec string) ([]db.Chat, error) {
	var chats []db.Chat
	tx := db.DBconn.Find(&chats, "Reciever= ?", rec)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return chats, nil
		}
		return nil, tx.Error
	}
	return chats, nil

}
