protos: 
	protoc --go_out=./ --go-grpc_out=./ ./protos/chat.proto
dev: 
	docker-compose up --build --remove-orphans 
client: 
	./client.exe 
 

