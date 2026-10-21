package main

import (
	"context"
	"errors"
	pb "handin3/gRPC"
	"log"
	"net"
	"strconv"
	"sync"
	"unicode/utf8"

	"google.golang.org/grpc"
)

type Client struct {
	Username string
	Messages chan *pb.ChatMessage
}

type ChitChatServer struct {
	pb.UnimplementedChitchatServer

	Clients map[string]*Client
	Clock   int32
	Mutex   sync.Mutex
}

func main() {
	server := &ChitChatServer{
		Clients: make(map[string]*Client),
	}

	server.startServer()
}

func (s *ChitChatServer) JoinBroadcast(
	user *pb.UserDetails,
	stream grpc.ServerStreamingServer[pb.ChatMessage],
) error {
	client := &Client{
		Username: user.Username,
		Messages: make(chan *pb.ChatMessage, 100),
	}

	s.Mutex.Lock()
	s.Clients[user.Username] = client
	s.Mutex.Unlock()

	log.Printf(
		"Server CLIENT_CONNECTED Client=%s",
		user.Username,
	)

	s.broadcast(
		user.Clock,
		user.Username,
		"Participant "+user.Username+" joined Chit Chat",
		true,
	)

	for {
		select {
		case message := <-client.Messages:
			err := stream.Send(message)
			if err != nil {
				s.removeClient(user.Username)
				return nil
			}

			log.Printf(
				"Server MESSAGE_DELIVERED Client=%s LogicalTime=%d Message=%s",
				user.Username,
				message.Clock,
				message.Msg,
			)

		case <-stream.Context().Done():
			s.removeClient(user.Username)
			return nil
		}
	}
}

func (s *ChitChatServer) Post(ctx context.Context, message *pb.ChatMessage) (*pb.Acknowledgement, error) {
	if !utf8.ValidString(message.Msg) {
		return &pb.Acknowledgement{}, errors.New("invalid message string")
	}

	if utf8.RuneCountInString(message.Msg) > 128 {
		return &pb.Acknowledgement{}, errors.New("message too long")
	}

	s.broadcast(
		message.Clock,
		message.Username,
		message.Msg,
		false,
	)

	return &pb.Acknowledgement{}, nil
}

func (s *ChitChatServer) broadcast(
	receivedClock int32,
	username string,
	text string,
	includeTimeInMessage bool,
) {
	s.Mutex.Lock()

	if receivedClock > s.Clock {
		s.Clock = receivedClock
	}

	s.Clock++
	logicalTime := s.Clock

	if includeTimeInMessage {
		text += " at logical time " +
			strconv.FormatInt(int64(logicalTime), 10)
	}

	message := &pb.ChatMessage{
		Clock:    logicalTime,
		Username: username,
		Msg:      text,
	}

	clients := make([]*Client, 0, len(s.Clients))

	for _, client := range s.Clients {
		clients = append(clients, client)
	}

	s.Mutex.Unlock()

	log.Printf(
		"Server MESSAGE_BROADCAST Client=%s LogicalTime=%d Message=%s",
		username,
		logicalTime,
		text,
	)

	for _, client := range clients {
		client.Messages <- message
	}
}

func (s *ChitChatServer) removeClient(username string) {
	s.Mutex.Lock()

	if _, exists := s.Clients[username]; !exists {
		s.Mutex.Unlock()
		return
	}

	delete(s.Clients, username)
	currentClock := s.Clock

	s.Mutex.Unlock()

	log.Printf(
		"Server CLIENT_DISCONNECTED Client=%s",
		username,
	)

	s.broadcast(
		currentClock,
		username,
		"Participant "+username+" left Chit Chat",
		true,
	)
}

func (s *ChitChatServer) startServer() {
	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterChitchatServer(grpcServer, s)

	log.Println("Server SERVER_STARTED Port=5050")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
