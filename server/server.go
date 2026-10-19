package server

import (
	proto "handin3/gRPC"
	"log"
	"net"

	"google.golang.org/grpc"
)

type ChitChatServer struct {
	proto.UnimplementedChitchatServer
	//needs more var's fill in later
}

func main() {
	server := &ChitChatServer{ //please fill this out later when more var's are added
	}
	server.start_server()
}

func (s *ChitChatServer) start_server() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("%s", "net.listen failed:"+err.Error())
	}

	proto.RegisterChitchatServer(grpcServer, s)

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalf("%s", "net.listen failed:"+err.Error())
	}

}
