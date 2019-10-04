package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"

	"github.com/hugozhu/grpc-demo/hello"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile           = flag.String("cert_file", "", "The TLS cert file")
	keyFile            = flag.String("key_file", "", "The TLS key file")
	port               = flag.Int("port", 10000, "The server port")
	isClient           = flag.Bool("client", false, "Running in client mode if true")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

func runClient() {
	fmt.Println("Hello from Client ...")
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer conn.Close()
	client := hello.NewHelloServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	req := &hello.HelloRequest{}
	req.Greeting = "Hello"
	resp, err := client.SayHello(ctx, req)
	if err != nil {
		cancel()
		fmt.Println(err)
	} else {
		fmt.Println(resp.GetReply())
	}
}

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
func (*HelloServiceServerImpl) LotsOfReplies(*hello.HelloRequest, hello.HelloService_LotsOfRepliesServer) error {
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

func runServer() {
	fmt.Println("Hello from Server ...")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = testdata.Path("server1.pem")
		}
		if *keyFile == "" {
			*keyFile = testdata.Path("server1.key")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	grpcServer := grpc.NewServer(opts...)
	hello.RegisterHelloServiceServer(grpcServer, &HelloServiceServerImpl{})
	grpcServer.Serve(lis)
}

func main() {
	flag.Parse()

	if *isClient {
		runClient()
	} else {
		runServer()
	}
}
