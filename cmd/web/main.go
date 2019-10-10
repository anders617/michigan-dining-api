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
var mockItems pb.Items

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
// Handler for GetItems request
//
func (s *server) GetItems(ctx context.Context, req *pb.ItemsRequest) (*pb.ItemsReply, error) {
	glog.Infof("GetItems req{%v}", req)
	return &pb.ItemsReply{Items: &mockItems}, nil
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

	// Register Server
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
}

func readProtoFromFile(path string, p proto.Message) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		glog.Fatalf("Failed to read in dininghalls text proto, %v", err)
	}
	proto.UnmarshalText(string(data), p)
}

func main() {
	flag.Parse()
	wg.Add(2)

	readProtoFromFile("cmd/web/dininghalls.proto.txt", &mockDiningHalls)
	readProtoFromFile("cmd/web/items.proto.txt", &mockItems)

	go serveGRPC()
	go serveHTTP()
	wg.Wait()
}
