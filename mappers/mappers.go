// TODO improve this
package mappers

import (
	"github.com/epicadk/grpc-chat/db"
	"github.com/epicadk/grpc-chat/models"
)

func ChatDBToProto(chat *db.Chat) *models.Message {
	return &models.Message{
		Id:       chat.ID,
		Sender:   chat.Sender,
		Body:     chat.Body,
		Receiver: chat.Receiver,
		Sent:     int64(chat.Sent),
	}
}

func ChatProtoToDB(msg *models.Message) *db.Chat {
	return &db.Chat{
		ID:       msg.Id,
		Sender:   msg.Sender,
		Receiver: msg.Receiver,
		Body:     msg.Body,
		Sent:     uint64(msg.Sent),
	}
}

func UserDBToProto(u *db.User) *models.User {
	return &models.User{
		UserID:      u.ID,
		Phonenumber: u.Phonenumber,
		DisplayName: u.DisplayName,
		Password:    u.Password,
	}
}

// warning does not copy userID
func UserProtoToDB(u *models.User) *db.User {
	return &db.User{
		ID:          u.UserID,
		Phonenumber: u.Phonenumber,
		DisplayName: u.DisplayName,
		Password:    u.Password,
	}
}
