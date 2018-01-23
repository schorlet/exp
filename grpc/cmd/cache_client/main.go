package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/schorlet/exp/grpc/interceptor"
	"github.com/schorlet/exp/grpc/rpc"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
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
	conn, err := grpc.Dial("localhost:5051",
		grpc.WithTransportCredentials(tlsCreds),
		interceptor.WithClientInterceptor(),
	)
	if err != nil {
		return fmt.Errorf("failed to dial server: %v", err)
	}
	cache := rpc.NewCacheClient(conn)

	// store
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	_, err = cache.Store(ctx, &rpc.StoreReq{
		AccountToken: "token",
		Key:          "gopher",
		Val:          []byte("con"),
	})
	if err != nil {
		log.Fatalf("Failed to store: %v", err)
	}

	// store (dry run)
	ctx, _ = context.WithTimeout(context.Background(), time.Second)
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("dry-run", "1"))
	_, err = cache.Store(ctx, &rpc.StoreReq{
		AccountToken: "token",
		Key:          "con",
		Val:          []byte("2017"),
	})
	if err != nil {
		log.Fatalf("Failed to store: %v", err)
	}

	// get
	ctx, _ = context.WithTimeout(context.Background(), 50*time.Millisecond)
	resp, err := cache.Get(ctx, &rpc.GetReq{Key: "gopher"})
	if err != nil {
		log.Fatalf("Failed to get: %v", err)
	}
	fmt.Printf("Got cached value: %s\n", resp.Val)

	// get, expects not found
	ctx, _ = context.WithTimeout(context.Background(), 50*time.Millisecond)
	resp, err = cache.Get(ctx, &rpc.GetReq{Key: "con"})
	if err == nil {
		log.Fatalf("Got cached value: %s\n", resp.Val)
	}
	if _, ok := status.FromError(err); !ok {
		log.Fatalf("Got unknown error: %v\n", err)
	}
	fmt.Printf("Got expected error: %v\n", err)

	return nil
}
