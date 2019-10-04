package client

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/hugozhu/grpc-demo/hello"
)

func SayHello(client hello.HelloServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	req := &hello.HelloRequest{}
	req.Greeting = "Hello " + time.Now().String()
	resp, err := client.SayHello(ctx, req)
	if err != nil {
		cancel()
		log.Println(err)
	} else {
		log.Println(resp.GetReply())
	}
}

func SayHelloInStream(client hello.HelloServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	req := &hello.HelloRequest{}
	req.Greeting = "Hello " + time.Now().String()
	stream, err := client.LotsOfReplies(ctx, req)
	if err != nil {
		cancel()
		return
	}
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.LotsOfReplies(_) = _, %v", client, err)
		}
		log.Println(reply.GetReply())
	}
}
