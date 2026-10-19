package server

import (
	"context"
	pb "handin3/gRPC"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

type ChitChatServer struct {
	pb.UnimplementedChitchatServer
	Time       []int32
	Component  []string
	EventType  []string
	Identifier []int32
}

func main() {
	server := &ChitChatServer{ //please fill this out later when more var's are added
	}
	server.start_server()
}

func (s *ChitChatServer) Broadcast(ud pb.UserDetails, sss grpc.ServerStreamingServer[pb.ChatMessage]) error {
	sss.SendMsg(pb.ChatMessage{Clock: 1, Msg: "Hello"})
	return nil
}

func (s *ChitChatServer) Post(ctx context.Context, cmsg *pb.ChatMessage) (*pb.Empty, error) {
	response := pb.Empty{}
	return &response, nil
}

func (s *ChitChatServer) start_server() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("%s", "net.listen failed:"+err.Error())
	}

	pb.RegisterChitchatServer(grpcServer, s)

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalf("%s", "net.listen failed:"+err.Error())
	}

}

func (*ChitChatServer) log(TIME int32, COMPONENT string, EVENTTYPE string, IDENTIFIER int32) {
	log.SetFlags(0)
	log.Println(strconv.FormatInt(int64(TIME), 10) + " " + COMPONENT + " " + strconv.FormatInt(int64(IDENTIFIER), 10) + " " + EVENTTYPE)
}
