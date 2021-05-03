package server

import (
	"context"
	"log"
	"sync"

	"github.com/epicadk/grpc-chat/dao"
	"github.com/epicadk/grpc-chat/models"
)

//Connection Represents a connection to the server.
type Connection struct {
	stream models.ChatService_LoginServer
	err    chan error
}

type Server struct {
	// Map of active connections
	Connections map[string]*Connection
}

var chatDao dao.ChatDao
var userDao dao.UserDao

func (s *Server) Login(loginRequest *models.LoginRequest, stream models.ChatService_LoginServer) error {
	err := userDao.CheckCredentials(loginRequest.Phonenumber, loginRequest.Password)
	log.Println(loginRequest.Phonenumber)

	if err != nil {
		return err
	}

	conn := &Connection{
		stream: stream,
		err:    make(chan error),
	}

	// TODO use userID instead of username to find the chat
	messages, err := chatDao.FindChat(loginRequest.Phonenumber)
	if err != nil {
		return err
	}

	for _, v := range messages {

		if err := conn.stream.Send(v); err != nil {
			return err
		}

		if err := chatDao.DeleteChat(v); err != nil {
			log.Fatal(err)
		}

	}
	s.Connections[loginRequest.Phonenumber] = conn

	// return is blocked till conn.err gets an error
	return <-conn.err
}

func (s *Server) SendChat(ctx context.Context, message *models.Message) (*models.Success, error) {
	log.Println(message)
	wg := sync.WaitGroup{}
	var f bool

	for k, con := range s.Connections {
		// can add multiple Recivers
		// if receiver is not here store in database
		if message.To == k {
			f = true
			wg.Add(1)

			go func(msg *models.Message, conn *Connection, wg *sync.WaitGroup) {
				defer wg.Done()

				if err := conn.stream.Send(msg); err != nil {
					conn.err <- err
					delete(s.Connections, k)
					f = false
				}
			}(message, con, &wg)
		}
	}

	wg.Wait()
	if !f {
		err := chatDao.CreateChat(message)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &models.Success{}, nil
}

// Probably create a different server for user related operations
func (s *Server) Register(ctx context.Context, user *models.User) (*models.Success, error) {
	err := userDao.Create(user)

	if err != nil {
		return &models.Success{Value: false}, err
	}

	return &models.Success{Value: true}, nil
}
