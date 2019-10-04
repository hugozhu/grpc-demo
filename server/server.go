package server

import (
	"context"
	"fmt"
	"io"
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
func (*HelloServiceServerImpl) LotsOfGreetings(stream hello.HelloService_LotsOfGreetingsServer) error {
	count := 0
	var greeting string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			resp := &hello.HelloResponse{
				Reply: greeting + " <= " + fmt.Sprintf("%d", count),
			}
			return stream.SendAndClose(resp)
		}
		if err != nil {
			return err
		}
		greeting = req.GetGreeting()
		count++
	}
}

//BidiHello is the implementation of interface function
func (*HelloServiceServerImpl) BidiHello(hello.HelloService_BidiHelloServer) error {
	return nil
}
