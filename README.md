#GRPC Demo

##generate grpc client/server code
protoc -I hello hello.proto --go_out=plugins=grpc:hello 
