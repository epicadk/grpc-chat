package server

import (
	"context"
	"log"
	"sync"

	"github.com/epicadk/grpc-chat/db"
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
	var messages []db.Chat
	res := db.DBconn.Where("Reciever = ?", loginRequset.Username).Find(&messages)
	if res.Error != nil {
		return res.Error
	}
	for _, v := range messages {
		// TODO probably create a DAO struct
		conn.stream.Send(&models.Message{
			Sender:   v.Sender,
			Body:     v.Body,
			Reciever: v.Reciever,
			Sent:     int64(v.Sent),
		})
	}

	db.DBconn.Delete(messages, "Reciever = ?", loginRequset.Username)
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
		res := db.DBconn.Create(&db.Chat{
			Sender:   message.Sender,
			Body:     message.Body,
			Reciever: message.Reciever,
			Sent:     uint64(message.Sent),
		})
		if res.Error != nil {
			log.Fatal(res.Error)
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
