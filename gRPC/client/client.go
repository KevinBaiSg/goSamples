package main

import (
	"context"
	"log"
	"os"

	pb "github.com/KevinBaiSg/goSamples/grpc/proto"
	"google.golang.org/grpc"
)

func main()  {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("connect server failed")
	}
	defer conn.Close()

	name := "hello"
	c := pb.NewGreeterClient(conn)

	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatal("SayHello failed")
	}

	log.Println(r.Message)
}
