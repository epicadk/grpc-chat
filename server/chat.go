package server

import (
	"context"
	"log"
	"sync"

	"github.com/epicadk/grpc-chat/models"
)

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
	s.Connections = append(s.Connections, conn)
	// return is blocked till conn.err gets an error
	return <-conn.err
}

func (s *Server) SendChat(ctx context.Context, message *models.Message) (*models.Success, error) {
	log.Println(message)
	wg := sync.WaitGroup{}
	for _, con := range s.Connections {
		// can add multiple Recivers
		// if reciever is not here store in database
		if message.Reciever == con.id {
			wg.Add(1)
			go handleMessages(message, con, &wg)
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
