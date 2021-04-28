package server

import (
	"context"
	"log"
	"sync"

	"github.com/epicadk/grpc-chat/db/dao"
	"github.com/epicadk/grpc-chat/models"
	"github.com/epicadk/grpc-chat/utils"
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

func (s *Server) Login(loginRequset *models.LoginRequest, stream models.ChatService_LoginServer) error {
	conn := &Connection{
		stream: stream,
		err:    make(chan error),
	}

	messages, err := chatDao.FindChat(loginRequset.Username)
	if err != nil {
		return err
	}

	for _, v := range messages {

		if err := conn.stream.Send(utils.ChatDbToProto(&v)); err != nil {
			return err
		}

		if err := chatDao.DeleteChat(&v); err != nil {
			log.Fatal(err)
		}

	}

	s.Connections[loginRequset.Username] = conn

	// return is blocked till conn.err gets an error
	return <-conn.err
}

func (s *Server) SendChat(ctx context.Context, message *models.Message) (*models.Success, error) {
	log.Println(message)

	wg := sync.WaitGroup{}
	var f bool

	for k, con := range s.Connections {
		// can add multiple Recivers
		// if reciever is not here store in database
		if message.Reciever == k {
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
		err := chatDao.SaveChat(message)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &models.Success{}, nil
}

func (s *Server) Register(ctx context.Context, user *models.User) (*models.RegisterResponse, error) {
	id, err := userDao.SaveUser(user)
	if err != nil {
		return nil, err
	}
	return &models.RegisterResponse{UserID: id}, nil
}
