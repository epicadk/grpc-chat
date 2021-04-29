package dao

import (
	"github.com/epicadk/grpc-chat/db"
	"github.com/epicadk/grpc-chat/mappers"
	"github.com/epicadk/grpc-chat/models"
)

type ChatDao struct{}

// TODO some validation
func (cd *ChatDao) CreateChat(msg *models.Message) error {
	return mappers.ChatProtoToDB(msg).SaveToDB()
}

func (cd *ChatDao) FindChat(userID string) ([]*models.Message, error) {
	chat := db.Chat{
		Reciever: userID,
	}
	chats, err := chat.FindChat()
	if err != nil {
		return nil, err
	}
	var msgs []*models.Message
	for _, v := range chats {
		msgs = append(msgs, mappers.ChatDBToProto(&v))
	}
	return msgs, nil
}

func (cd *ChatDao) DeleteChat(msg *models.Message) error {
	return mappers.ChatProtoToDB(msg).DeleteChat()
}
