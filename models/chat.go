package models

import "context"

type service struct {
}

var ActiveUsers []string
var messages map[string][]*Message

func (s *service) Login(loginRequest *LoginRequest, stream ChatService_LoginServer) error {
	ActiveUsers = append(ActiveUsers, loginRequest.Username)
	for _, v := range messages[loginRequest.Username] {
		if err := stream.Send(v); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) SendChat(c *context.Context, message *Message) *Success {
	messages[message.Reciever] = append(messages[message.Sender], message)
	return &Success{Value: true}
}
