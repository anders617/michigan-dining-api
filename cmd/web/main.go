package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"

	pb "github.com/MichiganDiningAPi/api/proto"
	"github.com/golang/protobuf/jsonpb"
	"github.com/google/logger"
	"google.golang.org/grpc"
)

const (
	grpcPort = ":1234"
	httpPort = ":1235"
)

var wg sync.WaitGroup

type server struct {
	pb.UnimplementedMDiningServer
}

func (s *server) GetDiningHalls(ctx context.Context, req *pb.DiningHallsRequest) (*pb.DiningHallsReply, error) {
	logger.Infof("GetDiningHalls ctx%v req%v", ctx, req)
	return &pb.DiningHallsReply{Test: "Hello, World!"}, nil
}

func serveGRPC() {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMDiningServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		logger.Fatalf("failed to server: %v", err)
	}
	wg.Done()
}

func serveHTTP() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("HTTP /")
		res, err := new(jsonpb.Marshaler).MarshalToString(&pb.DiningHallsReply{Test: "Hello, World!"})
		if err != nil {
			fmt.Fprintf(w, "Error")
		}
		fmt.Fprintf(w, res)
	})
	http.ListenAndServe(httpPort, nil)
	wg.Done()
}

func main() {
	logger.Init("Web", true, true, ioutil.Discard)
	defer logger.Close()
	wg.Add(2)
	go serveGRPC()
	go serveHTTP()
	wg.Wait()
}
