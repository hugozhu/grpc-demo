package client

import (
	"log"
	"time"
	"context"

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
