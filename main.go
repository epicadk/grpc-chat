package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/epicadk/grpc-chat/db"
	"github.com/epicadk/grpc-chat/models"
	"github.com/epicadk/grpc-chat/server"
	"github.com/epicadk/grpc-chat/utils"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("main: error loading .env file: %v", err)
	}

	db.SetupDB()

	Connections := make(map[string]*server.Connection)
	secret, check := os.LookupEnv("SECRET_KEY")
	if !check {
		log.Fatalf("cannot find secret key in environment")
	}
	jwtManager := utils.NewJWTManager(secret, 15*time.Minute)
	server := server.Server{Connections: Connections, JwtManager: jwtManager}
	go func() {
		for {
			for k := range Connections {
				log.Printf(k)
			}
			time.Sleep(time.Second * 300)
		}
	}()

	interceptor := utils.NewInterceptor(jwtManager, utils.Roles())
	grpcserver := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()), grpc.StreamInterceptor(interceptor.Stream()))

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	models.RegisterChatServiceServer(grpcserver, &server)
	if err := grpcserver.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
