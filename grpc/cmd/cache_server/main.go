package main

import (
	"log"
	"net"

	"github.com/schorlet/exp/grpc/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	if err := runServer(); err != nil {
		log.Fatalf("Failed to run cache server: %s\n", err)
	}
}

func runServer() error {
	tlsCreds, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	if err != nil {
		return err
	}
	srv := grpc.NewServer(grpc.Creds(tlsCreds))
	rpc.RegisterCacheServer(srv, &CacheService{
		store: make(map[string][]byte),
	})
	l, err := net.Listen("tcp", "localhost:5051")
	if err != nil {
		return err
	}
	// blocks until complete
	return srv.Serve(l)
}

// CacheService stores values in memory.
type CacheService struct {
	store map[string][]byte
}

// Get returns a value from the cache
func (s *CacheService) Get(ctx context.Context, req *rpc.GetReq) (*rpc.GetResp, error) {
	val, ok := s.store[req.Key]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Key not found %q", req.Key)
	}
	return &rpc.GetResp{Val: val}, nil
}

// Store sets a value into the cache
func (s *CacheService) Store(ctx context.Context, req *rpc.StoreReq) (*rpc.StoreResp, error) {
	s.store[req.Key] = req.Val
	return &rpc.StoreResp{}, nil
}
