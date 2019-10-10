package main

import (
	"context"
	"flag"
	"io/ioutil"
	"net"
	"net/http"
	"sync"

	pb "github.com/MichiganDiningAPI/api/proto"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

const (
	grpcPort = ":1234"
	httpPort = ":1235"
	grpcAddr = "localhost:1234"
)

var wg sync.WaitGroup
var mockDiningHalls pb.DiningHalls

type server struct {
}

//
// Handler for GetDiningHalls request
//
func (s *server) GetDiningHalls(ctx context.Context, req *pb.DiningHallsRequest) (*pb.DiningHallsReply, error) {
	glog.Infof("GetDiningHalls req{%v}", req)
	// Currently just returns static dining halls data that's checked into git
	return &pb.DiningHallsReply{DiningHalls: &mockDiningHalls}, nil
}

//
// Serves GRPC requests
//
func serveGRPC() {
	defer wg.Done()
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		glog.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMDiningServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		glog.Fatalf("failed to server: %v", err)
	}
}

//
// Proxies REST requests to GRPC server, converting to Proto
//
func serveHTTP() {
	defer wg.Done()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// Set the address to forward requests to to grpcAddr
	err := pb.RegisterMDiningHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		return
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	http.ListenAndServe(":8081", mux)

	wg.Done()
}

func main() {
	flag.Parse()
	wg.Add(2)
	data, err := ioutil.ReadFile("cmd/web/dininghalls.proto.txt")
	if err != nil {
		glog.Fatalf("Failed to read in dininghalls text proto, %v", err)
	}
	proto.UnmarshalText(string(data), &mockDiningHalls)
	go serveGRPC()
	go serveHTTP()
	wg.Wait()
}
