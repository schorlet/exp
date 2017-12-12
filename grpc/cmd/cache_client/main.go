package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/schorlet/exp/grpc/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	if err := runClient(); err != nil {
		log.Fatalf("Failed to run cache client: %s\n", err)
	}
}

func runClient() error {
	// connect
	// InsecureSkipVerify only for this example
	tlsCreds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	conn, err := grpc.Dial("localhost:5051", grpc.WithTransportCredentials(tlsCreds))
	if err != nil {
		return fmt.Errorf("failed to dial server: %v", err)
	}
	cache := rpc.NewCacheClient(conn)

	// store
	start := time.Now()
	_, err = cache.Store(context.Background(), &rpc.StoreReq{
		AccountToken: "token",
		Key:          "gopher",
		Val:          []byte("con"),
	})
	log.Printf("cache.Store duration %s", time.Since(start))
	if err != nil {
		log.Fatalf("failed to store: %v", err)
	}

	// get
	start = time.Now()
	resp, err := cache.Get(context.Background(), &rpc.GetReq{Key: "gopher"})
	log.Printf("cache.Get duration %s", time.Since(start))
	if err != nil {
		log.Fatalf("failed to get: %v", err)
	}
	fmt.Printf("Got cached value: %s\n", resp.Val)

	// get, expects not found
	resp, err = cache.Get(context.Background(), &rpc.GetReq{Key: "con"})
	if err == nil {
		log.Fatalf("Got cached value: %s\n", resp.Val)
	}
	if _, ok := status.FromError(err); !ok {
		log.Fatalf("Got unknown error: %v\n", err)
	}
	fmt.Printf("Got expected error: %v\n", err)

	return nil
}
