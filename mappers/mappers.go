// TODO improve this
package mappers

import (
	"github.com/epicadk/grpc-chat/db"
	"github.com/epicadk/grpc-chat/models"
)

func ChatDBToProto(chat *db.Chat) *models.Message {
	return &models.Message{
		Id:   chat.ID,
		From: chat.From,
		Body: chat.Body,
		To:   chat.To,
		Time: chat.Time,
	}
}

func ChatProtoToDB(msg *models.Message) *db.Chat {
	return &db.Chat{
		ID:   msg.Id,
		From: msg.From,
		To:   msg.To,
		Body: msg.Body,
		Time: msg.Time,
	}
}

func UserDBToProto(u *db.User) *models.User {
	return &models.User{
		Phonenumber: u.Phonenumber,
		DisplayName: u.DisplayName,
		Password:    u.Password,
	}
}

// warning does not copy userID
func UserProtoToDB(u *models.User) *db.User {
	return &db.User{
		Phonenumber: u.Phonenumber,
		DisplayName: u.DisplayName,
		Password:    u.Password,
	}
}
