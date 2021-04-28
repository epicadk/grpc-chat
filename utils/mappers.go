package utils

import (
	"github.com/epicadk/grpc-chat/db"
	"github.com/epicadk/grpc-chat/models"
)

func ChatDbToProto(chat *db.Chat) *models.Message {
	return &models.Message{
		Sender:   chat.Sender,
		Body:     chat.Body,
		Reciever: chat.Reciever,
		Sent:     int64(chat.Sent),
	}
}

func ChatProtoToD(msg *models.Message) *db.Chat {
	return &db.Chat{
		Sender:   msg.Sender,
		Reciever: msg.Reciever,
		Body:     msg.Body,
		Sent:     uint64(msg.Sent),
	}
}

func UserDbToProto(u *db.User) *models.User {
	return &models.User{
		UserID:      u.Id,
		DisplayName: u.Username,
		Password:    u.Password,
	}
}

// warning does not copy userID
func UserProtoToDb(u *models.User) *db.User {
	return &db.User{
		Username: u.DisplayName,
		Password: u.Password,
	}
}
