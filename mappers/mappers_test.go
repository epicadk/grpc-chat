package mappers_test

import (
	"log"
	"testing"

	"github.com/epicadk/grpc-chat/db"
	"github.com/epicadk/grpc-chat/mappers"
	"github.com/epicadk/grpc-chat/models"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestChatMappperDBtoProto(t *testing.T) {
	uuid, err := uuid.DefaultGenerator.NewV4()
	if err != nil {
		log.Fatal("Error Generating UUID", err)
	}
	testChat := db.Chat{
		ID:   uuid.String(),
		From: "from",
		Body: "Body",
		To:   "To",
		Time: 18181,
	}
	result := mappers.ChatDBToProto(&testChat)
	expected := models.Message{
		Id:   "",
		From: "from",
		Body: "Body",
		To:   "To",
		Time: 18181,
	}
	assert.Equal(t, expected, *result)
}

func TestChatMappperPrototoDB(t *testing.T) {
	uuid, err := uuid.DefaultGenerator.NewV4()
	if err != nil {
		log.Fatal("Error Generating UUID", err)
	}
	testChat := models.Message{
		Id:   uuid.String(),
		From: "from",
		Body: "Body",
		To:   "To",
		Time: 18181,
	}
	result := mappers.ChatProtoToDB(&testChat)
	expected := db.Chat{
		ID:   uuid.String(),
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
