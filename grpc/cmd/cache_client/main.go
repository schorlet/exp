package main

import (
	"fmt"
	"log"

	"github.com/schorlet/exp/grpc/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	if err := runClient(); err != nil {
		log.Fatalf("Failed to run cache client: %s\n", err)
	}
}

func runClient() error {
	// connect
	conn, err := grpc.Dial("localhost:5051", grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to dial server: %v", err)
	}
	cache := rpc.NewCacheClient(conn)

	// store
	_, err = cache.Store(context.Background(), &rpc.StoreReq{Key: "gopher", Val: []byte("con")})
	if err != nil {
		return fmt.Errorf("failed to store: %v", err)
	}

	// get
	resp, err := cache.Get(context.Background(), &rpc.GetReq{Key: "gopher"})
	if err != nil {
		return fmt.Errorf("failed to get: %v", err)
	}

	fmt.Printf("Got cached value: %s\n", resp.Val)
	return nil
}
