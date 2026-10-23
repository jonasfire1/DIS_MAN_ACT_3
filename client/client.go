package main

import (
	"bufio"
	"context"
	"fmt"
	pb "handin3/gRPC"
	"log"
	"os"
	"sync"
	"unicode/utf8"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ChitChatClient struct {
	Username string
	Clock    int32
	Mutex    sync.Mutex
	Client   pb.ChitchatClient
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Run the client with: go run ./client <username>")
	}

	conn, err := grpc.NewClient(
		"localhost:5050",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := &ChitChatClient{
		Username: os.Args[1],
		Client:   pb.NewChitchatClient(conn),
	}

	stream, err := client.Client.JoinBroadcast(
		context.Background(),
		&pb.UserDetails{
			Clock:    client.Clock,
			Username: client.Username,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	go client.receiveMessages(stream)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Write a message and press Enter:")

	for scanner.Scan() {
		message := scanner.Text()

		if message == "leave" {
			fmt.Println("Leaving the service.")
			client.Client.LeaveBroadcast(
				// yes this way leaves the potential for the client sending fake user details to disconnect other uses
				// but thats part of the fun of a chat service with friends isn't it?
				context.Background(),
				&pb.UserDetails{
					Clock:    client.Clock,
					Username: client.Username,
				},
			)
			break
		}

		if !utf8.ValidString(message) {
			fmt.Println("The message must be valid UTF-8.")
			continue
		}

		if utf8.RuneCountInString(message) > 128 {
			fmt.Println("The message cannot exceed 128 characters.")
			continue
		}

		client.sendMessage(message)
	}
}

func (c *ChitChatClient) sendMessage(text string) {
	c.Mutex.Lock()
	c.Clock++
	logicalTime := c.Clock
	c.Mutex.Unlock()

	ack, err := c.Client.Post(
		context.Background(),
		&pb.ChatMessage{
			Clock:    logicalTime,
			Username: c.Username,
			Msg:      text,
		},
	)
	if err != nil {
		log.Println(err)
		return
	}
	ack = ack // so go stops complaining

}

func (c *ChitChatClient) receiveMessages(
	stream grpc.ServerStreamingClient[pb.ChatMessage],
) {
	for {
		message, err := stream.Recv()
		if err != nil {
			return
		}

		c.Mutex.Lock()

		if message.Clock > c.Clock {
			c.Clock = message.Clock
		}

		c.Clock++

		c.Mutex.Unlock()

		log.Printf(
			"Client MESSAGE_RECEIVED Client=%s LogicalTime=%d Message=%s",
			c.Username,
			message.Clock,
			message.Msg,
		)
	}
}
