package main

import (
	"log"
	"net"
	"time"

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
		panic("error connecting to database")
	}

	err = db.DBconn.AutoMigrate(&db.Chat{})

	if err != nil {
		panic("error in auto migration")
	}
}

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
