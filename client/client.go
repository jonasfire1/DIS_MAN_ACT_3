package client

import (
	"log"

	proto "handin3/gRPC"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Connection failure")
	}
	client := proto.NewChitchatClient(conn)
}
