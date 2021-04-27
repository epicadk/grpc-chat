GenerateProtos: 
	protoc --go_out=plugins=grpc:./models ./protos/chat.proto
Dev: 
	docker-compose up --build --remove-orphans 
Client: 
	./client.exe 
 

