package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	tutorial "github.com/schorlet/exp/grpc-tutorial"

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
		tutorial.WithClientInterceptor(),
	)
	if err != nil {
		return fmt.Errorf("failed to dial server: %v", err)
	}
	cache := tutorial.NewCacheClient(conn)

	// store
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	_, err = cache.Store(ctx, &tutorial.StoreReq{
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
	_, err = cache.Store(ctx, &tutorial.StoreReq{
		AccountToken: "token",
		Key:          "con",
		Val:          []byte("2017"),
	})
	if err != nil {
		log.Fatalf("Failed to store: %v", err)
	}

	// get
	ctx, _ = context.WithTimeout(context.Background(), 50*time.Millisecond)
	resp, err := cache.Get(ctx, &tutorial.GetReq{Key: "gopher"})
	if err != nil {
		log.Fatalf("Failed to get: %v", err)
	}
	fmt.Printf("Got cached value: %s\n", resp.Val)

	// get, expects not found
	ctx, _ = context.WithTimeout(context.Background(), 50*time.Millisecond)
	resp, err = cache.Get(ctx, &tutorial.GetReq{Key: "con"})
	if err == nil {
		log.Fatalf("Got cached value: %s\n", resp.Val)
	}
	if _, ok := status.FromError(err); !ok {
		log.Fatalf("Got unknown error: %v\n", err)
	}
	fmt.Printf("Got expected error: %v\n", err)

	// dump
	ctx, _ = context.WithTimeout(context.Background(), 50*time.Millisecond)
	dump, err := cache.Dump(ctx, &tutorial.DumpReq{})
	if err != nil {
		log.Fatalf("Failed to dump: %v", err)
	}
	for _, item := range dump.Items {
		fmt.Printf("Got dump item: %s:%s\n", item.Key, item.Val)
	}

	return nil
}
