/*
packa
*/
package mappers

import (
	"github.com/epicadk/grpc-chat/db"
	"github.com/epicadk/grpc-chat/models"
)

// ChatDBToProto converts db.Chat struct to a models.Message struct
func ChatDBToProto(chat *db.Chat) *models.Message {
	return &models.Message{
		Id:   chat.ID,
		From: chat.From.Phonenumber,
		Body: chat.Body,
		To:   chat.To.Phonenumber,
		Time: chat.Time,
	}
}

// ChatDBToProto converts models.Message struct to a db.Chat struct
func ChatProtoToDB(msg *models.Message) *db.Chat {
	return &db.Chat{
		ID:   msg.Id,
		From: db.User{Phonenumber: msg.From},
		To:   db.User{Phonenumber: msg.To},
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

func UserProtoToDB(u *models.User) *db.User {
	return &db.User{
		Phonenumber: u.Phonenumber,
		DisplayName: u.DisplayName,
		Password:    u.Password,
	}
}
