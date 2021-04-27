package server

import (
	"context"
	"log"
	"sync"

	"github.com/epicadk/grpc-chat/db/dao"
	"github.com/epicadk/grpc-chat/models"
)

var chatDao dao.ChatDao

//Connection Represents a connection to the server.
type Connection struct {
	stream models.ChatService_LoginServer
	id     string
	active bool
	err    chan error
}

type Server struct {
	Connections []*Connection
}

func (s *Server) Login(loginRequset *models.LoginRequest, stream models.ChatService_LoginServer) error {
	conn := &Connection{
		stream: stream,
		id:     loginRequset.Username,
		active: true,
		err:    make(chan error),
	}
	messages, err := chatDao.FindChat(loginRequset.Username)
	if err != nil {
		return err
	}
	for _, v := range messages {
		// TODO probably create a DAO struct
		conn.stream.Send(&models.Message{
			Sender:   v.Sender,
			Body:     v.Body,
			Reciever: v.Reciever,
			Sent:     int64(v.Sent),
		})
		chatDao.DeleteChat(&v)

	}

	s.Connections = append(s.Connections, conn)
	// return is blocked till conn.err gets an error
	return <-conn.err
}

func (s *Server) SendChat(ctx context.Context, message *models.Message) (*models.Success, error) {
	log.Println(message)
	wg := sync.WaitGroup{}
	var f bool
	for _, con := range s.Connections {
		// can add multiple Recivers
		// if reciever is not here store in database
		if message.Reciever == con.id {
			f = true
			wg.Add(1)
			go handleMessages(message, con, &wg)
		}
	}
	if !f {
		err := chatDao.SaveChat(message)
		if err != nil {
			log.Fatal(err)
		}
	}

	wg.Wait()
	return &models.Success{}, nil
}

func handleMessages(msg *models.Message, conn *Connection, wg *sync.WaitGroup) {
	defer wg.Done()
	if conn.active {
		//ToDo remove connections
		if err := conn.stream.Send(msg); err != nil {
			conn.active = false
			conn.err <- err
		}
	}
}
