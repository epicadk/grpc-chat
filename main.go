package main

import (
	"log"
	"net"
	"time"

	"github.com/epicadk/grpc-chat/models"
	"github.com/epicadk/grpc-chat/server"
	"google.golang.org/grpc"
)

func main() {
	Connections := make(map[string]*server.Connection)
	server := server.Server{Connections: Connections}
	go func() {
		for {
			for k := range Connections {
				log.Printf(k)
			}
			time.Sleep(time.Second * 300)
		}
	}()
	grpcserver := grpc.NewServer()
	lis, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal(err)
	}

	models.RegisterChatServiceServer(grpcserver, &server)
	if err := grpcserver.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
