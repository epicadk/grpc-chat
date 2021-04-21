package main

import (
	"log"
	"net"

	"github.com/epicadk/grpc-chat/models"
	"github.com/epicadk/grpc-chat/server"
	"google.golang.org/grpc"
)

func main() {
	var Connections []*server.Connection
	server := &server.Server{Connections: Connections}

	grpcserver := grpc.NewServer()
	lis, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal(err)
	}

	models.RegisterChatServiceServer(grpcserver, server)
	grpcserver.Serve(lis)
}
