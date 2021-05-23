package mappers_test

import (
	"testing"

	"github.com/epicadk/grpc-chat/db"
	"github.com/epicadk/grpc-chat/mappers"
	"github.com/epicadk/grpc-chat/models"
	"github.com/stretchr/testify/assert"
)

func TestChatMappperDBtoProto(t *testing.T) {
	testChat := db.Chat{
		ID:   99090,
		From: "from",
		Body: "Body",
		To:   "To",
		Time: 18181,
	}
	result := mappers.ChatDBToProto(&testChat)
	expected := models.Message{
		Id:   99090,
		From: "from",
		Body: "Body",
		To:   "To",
		Time: 18181,
	}
	assert.Equal(t, expected, *result)
}

func TestChatMappperPrototoDB(t *testing.T) {
	testChat := models.Message{
		Id:   99090,
		From: "from",
		Body: "Body",
		To:   "To",
		Time: 18181,
	}
	result := mappers.ChatProtoToDB(&testChat)
	expected := db.Chat{
		ID:   99090,
		From: "from",
		Body: "Body",
		To:   "To",
		Time: 18181,
	}
	assert.Equal(t, expected, *result)
}

func TestUserMappperDBtoProto(t *testing.T) {
	testUser := db.User{
		Phonenumber: "9999999999",
		DisplayName: "Random Display Name",
		Password:    "mPass",
	}
	result := mappers.UserDBToProto(&testUser)
	expected := models.User{
		Phonenumber: "9999999999",
		DisplayName: "Random Display Name",
		Password:    "mPass",
	}
	assert.Equal(t, expected, *result)
}

func TestUserMappperProtoToDB(t *testing.T) {
	testUser := models.User{
		Phonenumber: "9999999999",
		DisplayName: "Random Display Name",
		Password:    "mPass",
	}
	result := mappers.UserProtoToDB(&testUser)
	expected := db.User{
		Phonenumber: "9999999999",
		DisplayName: "Random Display Name",
		Password:    "mPass",
	}
	assert.Equal(t, expected, *result)
}
