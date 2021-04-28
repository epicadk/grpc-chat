package dao

import (
	"github.com/epicadk/grpc-chat/db"
	"github.com/epicadk/grpc-chat/models"
	"github.com/epicadk/grpc-chat/utils"
)

type ChatDao struct{}

// TODO some validation
func (cd *ChatDao) CreateChat(msg *models.Message) error {
	return utils.ChatProtoToD(msg).SaveToDB()
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
		msgs = append(msgs, utils.ChatDbToProto(&v))
	}
	return msgs, nil
}

func (cd *ChatDao) DeleteChat(msg *models.Message) error {
	return utils.ChatProtoToD(msg).DeleteChat()
}
