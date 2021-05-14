package server

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/epicadk/grpc-chat/dao"
	"github.com/epicadk/grpc-chat/models"
)

// Connection Represents a connection to the server.
type Connection struct {
	stream models.ChatService_ConnectServer // the stream of the user
	err    chan error                       // the channel for the error
}

type Server struct {
	Connections map[string]*Connection // Map of active connections
	models.UnimplementedChatServiceServer
}

var (
	chatDao dao.ChatDao
	userDao dao.UserDao
)

func (s *Server) Login(ctx context.Context, loginRequest *models.LoginRequest) (*models.Success, error) {
	_, err := userDao.CheckCredentials(loginRequest.Phonenumber, loginRequest.Password)

	if err != nil {
		return &models.Success{Value: false}, err
	}

	return &models.Success{Value: true}, nil
}

func (s *Server) Connect(phone *models.Phone, stream models.ChatService_ConnectServer) error {
	conn := &Connection{
		stream: stream,
		err:    make(chan error),
	}

	messages, err := chatDao.FindChat(phone.Phonenumber)
	if err != nil {
		return err
	}

	s.Connections[phone.Phonenumber] = conn

	for _, v := range messages {
		go func(message *models.Message) {
			if err := conn.stream.Send(message); err != nil {
				conn.err <- err
			}

			// Delete only after message has been sent
			if err := chatDao.DeleteChat(message); err != nil {
				log.Fatal(err)
			}

		}(v)
	}

	// return is blocked till conn.err gets an error
	return <-conn.err
}

func (s *Server) SendChat(ctx context.Context, message *models.Message) (*models.Success, error) {
	log.Println(message)
	wg := sync.WaitGroup{}
	to, ok := s.Connections[message.To]

	timeMilli := time.Now().UnixNano() / 1e6
	message.Time = uint64(timeMilli)
	// can add multiple Recivers
	// if receiver is not here store in database
	if ok {
		wg.Add(1)

		go func(msg *models.Message, conn *Connection, wg *sync.WaitGroup) {
			defer wg.Done()

			if err := conn.stream.Send(msg); err != nil {
				conn.err <- err
				delete(s.Connections, message.To)
			}
		}(message, to, &wg)
	}
	wg.Wait()
	if !ok {
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
