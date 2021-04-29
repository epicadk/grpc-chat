protos: 
	protoc --go_out=plugins=grpc:./ ./protos/chat.proto
dev: 
	docker-compose up --build --remove-orphans 
client: 
	./client.exe 
 

