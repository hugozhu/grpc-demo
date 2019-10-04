package server

import (
	"context"
	"log"
	"sync"

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
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		resp := &hello.HelloResponse{
			Reply: req.GetGreeting() + " <= 1",
		}
		if err := stream.Send(resp); err != nil {
			log.Println(err)
		}
		wg.Done()
	}()

	go func() {
		resp := &hello.HelloResponse{
			Reply: req.GetGreeting() + " <= 2",
		}
		if err := stream.Send(resp); err != nil {
			log.Println(err)
		}
		wg.Done()
	}()
	wg.Wait()
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
