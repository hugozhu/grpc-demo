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

func SayHelloWithStreamOut(client hello.HelloServiceClient) {
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

func SayHelloWithStreamIn(client hello.HelloServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	req := &hello.HelloRequest{}
	req.Greeting = "Hello " + time.Now().String()
	stream, err := client.LotsOfGreetings(ctx)
	if err != nil {
		cancel()
		return
	}
	for i := 0; i < 3; i++ {
		if err := stream.Send(req); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("%v.Send(%v) = %v", stream, req, err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Println(reply.GetReply())
}

func SayHelloWithStreamBidi(client hello.HelloServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	req := &hello.HelloRequest{}
	req.Greeting = "Hello " + time.Now().String()
	stream, err := client.BidiHello(ctx)
	if err != nil {
		cancel()
		return
	}
	waitc := make(chan struct{})

	//read
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			log.Println(resp.GetReply())
		}
	}()

	//write
	for i := 0; i < 3; i++ {
		if err := stream.Send(req); err != nil {
			log.Fatalf("%v.Send(%v) = %v", stream, req, err)
		}
	}
	stream.CloseSend()
	<-waitc
}
