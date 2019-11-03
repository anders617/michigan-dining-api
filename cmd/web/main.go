package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	pb "github.com/MichiganDiningAPI/api/proto"
	"github.com/MichiganDiningAPI/internal/web/mdiningserver"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

//
// Launches a mdining-api server that handles http/rest, grpc, and grpc-web requests
//

const proxiedGrpcPort = "3000"

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
func serveGRPC(port string, server *mdiningserver.Server) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		glog.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Register Server
	pb.RegisterMDiningServer(s, server)

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

	mDiningServer := mdiningserver.New()

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
	pb.RegisterMDiningServer(grpcS, mDiningServer)

	// HTTP
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// Set the address to forward requests to to grpcAddr
	err = pb.RegisterMDiningHandlerFromEndpoint(ctx, mux, "localhost:"+proxiedGrpcPort, opts)
	grpcServer := grpc.NewServer()
	// Wrap it in a grpcweb handler in order to also serve grpc-web requests
	wrappedGrpc := grpcweb.WrapServer(grpcServer)
	grpcWebHandler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if wrappedGrpc.IsGrpcWebRequest(req) {
			wrappedGrpc.ServeHTTP(resp, req)
			return
		}
		// Fall back to other servers.
		mux.ServeHTTP(resp, req)
	})
	httpS := &http.Server{
		Handler: allowCORS(grpcWebHandler),
	}

	// Use the muxed listeners for your servers.
	// One GRPC server to handle proxied http requests
	go serveGRPC(proxiedGrpcPort, mDiningServer)
	// Second GRPC server to handle direct GRPC requests
	go grpcS.Serve(grpcL)
	// HTTP Server To Proxy Requests to First GRPC Server
	go httpS.Serve(httpL)

	// Start serving!
	m.Serve()
}
