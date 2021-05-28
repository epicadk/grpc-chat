package server

import (
	"context"
	"errors"
	"io"
	"log"
	"sync"
	"time"

	"github.com/epicadk/grpc-chat/dao"
	"github.com/epicadk/grpc-chat/models"
	"github.com/epicadk/grpc-chat/utils"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Connection Represents a connection to the server.
type Connection struct {
	stream models.ChatService_ConnectServer // the stream of the user
	err    chan error                       // the channel for the error
}

type Server struct {
	Connections map[string]*Connection // Map of active connections
	JwtManager  *utils.JWTManager      // JWT Manager
	models.UnimplementedChatServiceServer
}

var (
	chatDao dao.ChatDao
	userDao dao.UserDao
)

func (s *Server) Login(ctx context.Context, loginRequest *models.LoginRequest) (res *models.LoginResponse, err1 error) {
	res = &models.LoginResponse{}

	user, err := userDao.CheckCredentials(loginRequest.Phonenumber, loginRequest.Password)
	if err != nil {
		res.Status = &models.Success{Value: false}
		return res, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}

	if user == nil || utils.ComparePassword(user.Password, loginRequest.Password) != nil {
		res.Status = &models.Success{Value: false}
		return res, status.Errorf(codes.NotFound, "incorrect username/password")
	}

	token, err := s.JwtManager.Generate(user.Phonenumber)
	if err != nil {
		res.Status = &models.Success{Value: false}
		return res, status.Errorf(codes.Internal, "cannot generate access token")
	}

	res.Status = &models.Success{Value: true}
	res.AccessToken = token
	return res, nil
}

func (s *Server) Connect(stream models.ChatService_ConnectServer) error {
	userConnection := &Connection{
		stream: stream,
		err:    make(chan error),
	}
	md, ok := metadata.FromIncomingContext(stream.Context())
	// should not happen because interceptor
	if !ok {
		return errors.New("no metadata")
	}
	user := md.Get("user")[0]
	messages, err := chatDao.FindChat(user)
	if err != nil {
		return err
	}
	s.Connections[user] = userConnection
	log.Println(len(s.Connections), user)
	for _, v := range messages {
		go func(message *models.Message) {
			if err := userConnection.stream.Send(message); err != nil {
				userConnection.err <- err
			}

			// Delete only after message has been sent
			if err := chatDao.DeleteChat(message); err != nil {
				log.Fatal(err)
			}

		}(v)
	}
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				userConnection.err <- err
				return
			}
			wg := sync.WaitGroup{}
			timeMilli := time.Now().UnixNano() / 1e6
			msg.Time = uint64(timeMilli)
			uuid, err := uuid.DefaultGenerator.NewV4()
			msg.Id = uuid.String()
			log.Println(msg)
			to, ok := s.Connections[msg.To]
			if ok {
				wg.Add(1)
				go func(msg *models.Message, conn *Connection, wg *sync.WaitGroup) {
					defer wg.Done()
					if err := conn.stream.Send(msg); err != nil {
						conn.err <- err
						delete(s.Connections, msg.To)
					}
				}(msg, to, &wg)
			}
			wg.Wait()
			if !ok {
				err := chatDao.CreateChat(msg)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	// return is blocked till conn.err gets an error
	return <-userConnection.err
}

// Probably create a different server for user related operations
func (s *Server) Register(ctx context.Context, user *models.User) (*models.Success, error) {
	err := userDao.Create(user)

	if err != nil {
		return &models.Success{Value: false}, err
	}

	return &models.Success{Value: true}, nil
}
