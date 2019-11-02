package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	pb "github.com/MichiganDiningAPI/api/proto"
	"github.com/MichiganDiningAPI/db/dynamoclient"
	"github.com/MichiganDiningAPI/internal/processing/mdiningprocessing"
	"github.com/MichiganDiningAPI/internal/util/date"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

var wg sync.WaitGroup
var mockDiningHalls *pb.DiningHalls = &pb.DiningHalls{}
var mockItems *pb.Items = &pb.Items{}
var mockFilterableEntries *pb.FilterableEntries = &pb.FilterableEntries{}

const proxiedGrpcPort = "3000"

type server struct {
	dc *dynamoclient.DynamoClient
}

func NewServer() *server {
	server := server{dc: dynamoclient.New()}
	return &server
}

//
// Handler for GetDiningHalls request
//
func (s *server) GetDiningHalls(ctx context.Context, req *pb.DiningHallsRequest) (*pb.DiningHallsReply, error) {
	glog.Infof("GetDiningHalls req{%v}", req)
	// Currently just returns static dining halls data that's checked into git
	return &pb.DiningHallsReply{DiningHalls: mockDiningHalls.DiningHalls}, nil
}

//
// Handler for GetItems request
//
func (s *server) GetItems(ctx context.Context, req *pb.ItemsRequest) (*pb.ItemsReply, error) {
	glog.Infof("GetItems req{%v}", req)
	return &pb.ItemsReply{Items: mockItems.Items}, nil
}

func (s *server) GetFilterableEntries(ctx context.Context, req *pb.FilterableEntriesRequest) (*pb.FilterableEntriesReply, error) {
	glog.Infof("GetFilterableEntries req{%v}", req)
	return &pb.FilterableEntriesReply{FilterableEntries: mockFilterableEntries.FilterableEntries}, nil
}

func (s *server) GetAll(ctx context.Context, req *pb.AllRequest) (*pb.AllReply, error) {
	glog.Infof("GetAll req{%v}", req)
	return &pb.AllReply{DiningHalls: mockDiningHalls.DiningHalls, Items: mockItems.Items, FilterableEntries: mockFilterableEntries.FilterableEntries}, nil
}

func (s *server) GetMenu(ctx context.Context, req *pb.MenuRequest) (*pb.MenuReply, error) {
	glog.Infof("GetMenu req{%v}", req)
	diningHall, date, meal := &req.DiningHall, &req.Date, &req.Meal
	if *diningHall == "" {
		diningHall = nil
	}
	if *date == "" {
		date = nil
	}
	if *meal == "" {
		meal = nil
	}
	menus, err := s.dc.QueryMenus(diningHall, date, meal)
	if err != nil {
		glog.Infof("GetMenu Error %s", err)
		return nil, err
	}
	glog.Infof("GetMenu res{%d menus}", len(*menus))
	return &pb.MenuReply{Menus: *menus}, nil
}

func (s *server) GetFood(ctx context.Context, req *pb.FoodRequest) (*pb.FoodReply, error) {
	glog.Infof("GetFood req{%v}", req)
	name, date, startDate, endDate := &req.Name, &req.Date, &req.StartDate, &req.EndDate
	if *name == "" {
		name = nil
	}
	if *date == "" {
		date = nil
	}
	if *startDate == "" {
		startDate = nil
	}
	if *endDate == "" {
		endDate = nil
	}
	var foods *[]*pb.Food
	var err error
	if startDate != nil || endDate != nil {
		foods, err = s.dc.QueryFoodsDateRange(name, startDate, endDate)
	} else {
		foods, err = s.dc.QueryFoods(name, date)
	}
	if err != nil {
		glog.Infof("GetFood Error %s", err)
		return nil, err
	}
	glog.Infof("GetFood res{%d foods}", len(*foods))
	return &pb.FoodReply{Foods: *foods}, nil
}

// allowCORS allows Cross Origin Resoruce Sharing from any origin.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	})
}

//
// Serves GRPC requests
//
func serveGRPC(port string) {
	defer wg.Done()
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		glog.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Register Server
	pb.RegisterMDiningServer(s, NewServer())

	glog.Infof("Serving GRPC Requests on %s", port)
	if err := s.Serve(lis); err != nil {
		glog.Fatalf("failed to server: %v", err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	flag.Parse()
	wg.Add(2)

	dc := dynamoclient.New()
	var err error
	mockDiningHalls, err = dc.QueryDiningHalls()
	if err != nil {
		glog.Fatalf("QueryDiningHalls err %s", err)
	}
	glog.Infof("QueryDiningHalls Success")

	var foods *[]*pb.Food
	// Get all foods after today
	startDate := date.Format(date.DayStart(date.Now()))
	foods, err = dc.QueryFoodsDateRange(nil, &startDate, nil)
	if err != nil {
		glog.Fatalf("QueryFoodsDateRange err %s", err)
	}
	glog.Infof("QueryFoodsDateRange Success")
	mockItems = mdiningprocessing.FoodsToItems(foods)

	mockFilterableEntries = mdiningprocessing.ItemsToFilterableEntries(mockItems)

	// Create the main listener.
	glog.Infof("Listening on port " + port)
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	// Create a cmux.
	m := cmux.New(l)

	// Match connections in order:
	// First grpc, then HTTP, and otherwise Go RPC/TCP.
	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())

	// Create your protocol servers.
	grpcS := grpc.NewServer()

	// Register Server
	pb.RegisterMDiningServer(grpcS, NewServer())

	// HTTP
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// Set the address to forward requests to to grpcAddr
	err = pb.RegisterMDiningHandlerFromEndpoint(ctx, mux, "localhost:"+proxiedGrpcPort, opts)
	httpS := &http.Server{
		Handler: allowCORS(mux),
	}

	// Use the muxed listeners for your servers.
	// One GRPC server to handle proxied http requests
	go serveGRPC(proxiedGrpcPort)
	// Second GRPC server to handle direct GRPC requests
	go grpcS.Serve(grpcL)
	// HTTP Server To Proxy Requests to First GRPC Server
	go httpS.Serve(httpL)

	// Start serving!
	m.Serve()

	wg.Wait()
}
