generateProtos: 
	protoc --go_out=plugins=grpc:./models ./protos/chat.proto
dev: 
	docker-compose up --build --remove-orphans 
client: 
	./client.exe 
 

