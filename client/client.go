package main

// Simple client in go
import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/epicadk/grpc-chat/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ClientInterceptor struct {
	authMethods map[string]bool
	accessToken string
}

var client models.ChatServiceClient
var a string

func main() {
	interceptor, err := NewInterceptor(authMethods())

	if err != nil {
		log.Fatal("cannot create client interceptor: ", err)
	}

	conn, err := grpc.Dial(":8080", grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor.Unary()), grpc.WithStreamInterceptor(interceptor.Stream()))
	if err != nil {
		log.Fatal("cannot create GPRC client", err)
	}

	defer conn.Close()
	client = models.NewChatServiceClient(conn)

	fmt.Scanln(&a)

	var login string
	var password string
	fmt.Scanln(&login)
	fmt.Scanln(&password)

	switch a {
	case "r":
		sendRegister(login, password)
	case "l":
		err := sendLogin(&models.LoginRequest{Phonenumber: login, Password: password}, interceptor)
		if err != nil {
			log.Fatal(err)
		}
	}
	cStream, err := client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	waitc := make(chan struct{})
	// Get message from server
	go func() {
		for {
			msg, err := cStream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			log.Println(msg)
		}
	}()
	go func() {
		for {
			var body string
			var To string
			fmt.Scanln(&body)
			fmt.Scanln(&To)
			msg := &models.Message{
				From: login,
				Body: body,
				To:   To,
			}
			if err := cStream.Send(msg); err != nil {
				log.Fatalf("Failed to send a note: %v", err)
			}
		}
	}()
	<-waitc

}

func NewInterceptor(authMethods map[string]bool) (*ClientInterceptor, error) {
	interceptor := &ClientInterceptor{
		authMethods: authMethods,
	}

	return interceptor, nil
}

func (interceptor *ClientInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.Printf("--> unary interceptor: %s", method)

		if interceptor.authMethods[method] {
			return invoker(interceptor.attachToken(ctx), method, req, reply, cc, opts...)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (interceptor *ClientInterceptor) Stream() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		log.Printf("--> stream interceptor: %s", method)

		if interceptor.authMethods[method] {
			return streamer(interceptor.attachToken(ctx), desc, cc, method, opts...)
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}

func (interceptor *ClientInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", interceptor.accessToken)
}

func authMethods() map[string]bool {
	const servicePath = "/chat.ChatService/"

	return map[string]bool{
		servicePath + "Login":    false,
		servicePath + "Register": false,
		servicePath + "Connect":  true,
		servicePath + "SendChat": true,
	}
}

func sendLogin(req *models.LoginRequest, interceptor *ClientInterceptor) error {
	res, err := client.Login(context.Background(), req)
	if err != nil {
		return err
	}

	interceptor.accessToken = res.AccessToken
	return nil
}

func sendRegister(login, password string) {
	_, err := client.Register(context.Background(), &models.User{
		Phonenumber: login,
		Password:    password,
	})
	if err != nil {
		log.Fatal(err)
	}
}
