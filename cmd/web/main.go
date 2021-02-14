package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/MichiganDiningAPI/api/analytics/analyticsclient"
	"github.com/MichiganDiningAPI/internal/web/mdiningserver"
	pb "github.com/anders617/mdining-proto/proto/mdining"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

//
// Launches a mdining-api server that handles http/rest, grpc, and grpc-web requests
//

const proxiedGrpcPort = "5982"

var analytics *analyticsclient.AnalyticsClient = analyticsclient.New()

// preflightHandler adds the necessary headers in order to serve
// CORS from any origin using the methods "GET", "HEAD", "POST", "PUT", "DELETE"
// We insist, don't do this without consideration in production systems.
func preflightHandler(w http.ResponseWriter, r *http.Request) {
	// headers := []string{"Content-Type", "Accept", "Authorization", "x-user-agent"}
	w.Header().Set("Access-Control-Allow-Headers", "*") //strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	glog.Infof("preflight request for %s", r.URL.Path)
}

// allowCORS allows Cross Origin Resoruce Sharing from any origin.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		glog.Infof("serving http for %s", r.URL.Path)
		// Asynchronously send analytics
		go analytics.SendHit(r)
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

	// Read index.html and favicon.ico into memory
	indexHTML, e := ioutil.ReadFile("public/index.html")
	if e != nil {
		glog.Fatalf("Error reading public/index.html %s", e)
	}
	var favicon []byte
	favicon, e = ioutil.ReadFile("public/favicon.ico")
	if e != nil {
		glog.Fatalf("Error reading public/favicon.ico", e)
	}

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
	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())

	// HTTP
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// Set the address to forward requests to to grpcAddr
	err = pb.RegisterMDiningHandlerFromEndpoint(ctx, mux, "localhost:"+proxiedGrpcPort, opts)
	grpcServer := grpc.NewServer()
	// Register Server
	pb.RegisterMDiningServer(grpcServer, mDiningServer)
	// Wrap it in a grpcweb handler in order to also serve grpc-web requests
	wrappedGrpc := grpcweb.WrapServer(grpcServer, grpcweb.WithAllowedRequestHeaders([]string{"*"}))
	grpcWebHandler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if wrappedGrpc.IsGrpcWebRequest(req) {
			wrappedGrpc.ServeHTTP(resp, req)
			return
		}
		if req.URL.Path == "/" {
			resp.Write(indexHTML)
			return
		}
		if req.URL.Path == "/favicon.ico" {
			resp.Header().Set("Content-Type", "image/x-icon")
			resp.Write(favicon)
			return
		}
		if req.URL.Path == "/healthcheck" {
			if mDiningServer.IsAvailable() {
				resp.Write([]byte("OK"))
				return
			}
			http.Error(resp, "Unavailable", http.StatusInternalServerError)
			return
		}
		// Fall back to other servers.
		mux.ServeHTTP(resp, req)
	})
	httpS := &http.Server{
		Handler: allowCORS(grpcWebHandler),
	}

	// Create your protocol servers.
	grpcS := grpc.NewServer()

	// Register Server
	pb.RegisterMDiningServer(grpcS, mDiningServer)

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
