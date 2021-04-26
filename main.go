package main

import (
	"log"
	"net"

	"github.com/epicadk/grpc-chat/db"
	"github.com/epicadk/grpc-chat/models"
	"github.com/epicadk/grpc-chat/server"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	var err error
	// TODO use env vars
	var dns = "host=db user=postgres password=postgres dbname=chats port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	db.DBconn, err = gorm.Open(postgres.Open(dns), &gorm.Config{})

	if err != nil {
		panic("Error Connecting to database")
	}

	db.DBconn.AutoMigrate(&db.Chat{})
}

func main() {
	var Connections []*server.Connection
	server := server.Server{Connections: Connections}

	grpcserver := grpc.NewServer()
	lis, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal(err)
	}

	models.RegisterChatServiceServer(grpcserver, &server)
	grpcserver.Serve(lis)
}
