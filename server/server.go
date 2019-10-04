package server

import (
	"context"

	"github.com/hugozhu/grpc-demo/hello"
)

//HelloServiceServerImpl is the implementation of HelloServiceServer
type HelloServiceServerImpl struct {
}

//SayHello is the implementation of interface function
func (*HelloServiceServerImpl) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	resp := &hello.HelloResponse{
		Reply: req.GetGreeting(),
	}
	return resp, nil
}

//LotsOfReplies is the implementation of interface function
func (*HelloServiceServerImpl) LotsOfReplies(req *hello.HelloRequest, stream hello.HelloService_LotsOfRepliesServer) error {
	resp := &hello.HelloResponse{
		Reply: req.GetGreeting(),
	}
	if err := stream.Send(resp); err != nil {
		return err
	}
	if err := stream.Send(resp); err != nil {
		return err
	}
	return nil
}

//LotsOfGreetings is the implementation of interface function
func (*HelloServiceServerImpl) LotsOfGreetings(hello.HelloService_LotsOfGreetingsServer) error {
	return nil
}

//BidiHello is the implementation of interface function
func (*HelloServiceServerImpl) BidiHello(hello.HelloService_BidiHelloServer) error {
	return nil
}
