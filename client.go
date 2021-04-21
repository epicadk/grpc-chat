package main

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
	fmt.Println(a)
	sendMessage(&models.LoginRequest{Username: a, Password: "random", Active: true})
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			body := scanner.Text()
			scanner.Scan()
			msg := &models.Message{
				Sender:   a,
				Body:     body,
				Reciever: scanner.Text(),
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

func sendMessage(req *models.LoginRequest) error {
	stream, err := client.Login(context.Background(), req)
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
