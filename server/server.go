package server

import (
	proto "handin3/gRPC"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

type ChitChatServer struct {
	proto.UnimplementedChitchatServer
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

func (*ChitChatServer) log(TIME int32, COMPONENT string, EVENTTYPE string, IDENTIFIER int32) {
	log.SetFlags(0)
	log.Println(strconv.FormatInt(int64(TIME), 10) + " " + COMPONENT + " " + strconv.FormatInt(int64(IDENTIFIER), 10) + " " + EVENTTYPE)
}
