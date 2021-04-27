package dao

import (
	"github.com/epicadk/grpc-chat/db"
	"github.com/epicadk/grpc-chat/models"
	"github.com/epicadk/grpc-chat/utils/mappers"
)

type ChatDao struct{}

func (cd ChatDao) SaveChat(chat *models.Message) error {

	tx := db.DBconn.Create(mappers.ToDB(chat))
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
	tx := db.DBconn.Find(chats, "Reciever= ?", rec)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return chats, nil

}