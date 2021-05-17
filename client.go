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
)

var client models.ChatServiceClient
var wg *sync.WaitGroup
var a string

func init() {
	wg = &sync.WaitGroup{}
}

func main() {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
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
		success, _ := sendLogin(&models.LoginRequest{Phonenumber: login, Password: password})
		if success.Value == true {
			makeConnection(&models.Phone{Phonenumber: login})
		}
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

func sendLogin(req *models.LoginRequest) (*models.Success, error) {
	success, err := client.Login(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	return success, err
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
