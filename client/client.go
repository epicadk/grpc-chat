package main

// Simple client in go
import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/epicadk/grpc-chat/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ClientInterceptor struct {
	authMethods map[string]bool
	accessToken string
}

var client models.ChatServiceClient
var wg *sync.WaitGroup
var a string

func init() {
	wg = &sync.WaitGroup{}
}

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
		makeConnection(&models.Phone{Phonenumber: login})
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			body := scanner.Text()
			scanner.Scan()
			msg := &models.Message{
				From: login,
				Body: body,
				To:   scanner.Text(),
			}
			_, err := client.SendChat(context.Background(), msg)
			if err != nil {
				log.Fatal(err)
				break
			}
		}
	}()
	wg.Wait()
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

func makeConnection(phone *models.Phone) error {
	stream, err := client.Connect(context.Background(), phone)
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Fatal(err)
				break
			}
			log.Println(msg)
		}
	}()
	return err
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
