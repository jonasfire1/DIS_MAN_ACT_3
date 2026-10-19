package main

import (
	"context"
	pb "handin3/gRPC"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Connection failure")
	}
	client := pb.NewChitchatClient(conn)

	ssc, err := client.Broadcast(context.Background(), &pb.UserDetails{})
	log.Println("printing")
	for {
		msg, err := ssc.Recv()
		if err != nil {
			log.Fatalf(err.Error())
			break
		}
		log.Println(msg.Msg)
		// msg is was recieved and is supposed to be handled here
	}
	log.Println("exit")
}
