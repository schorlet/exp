package tutorial_test

import (
	"testing"
	"time"

	"github.com/schorlet/exp/grpc-tutorial"
	"github.com/schorlet/exp/grpc-tutorial/api"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func withClient(t *testing.T, fn func(api.CacheClient)) {
	var (
		adress = "localhost:5051"
		cert   = "cert.pem"
		key    = "key.pem"
	)

	// start a new server
	stop, err := tutorial.Serve(adress, cert, key)
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer stop()

	// create cache client
	tlsCreds, err := credentials.NewClientTLSFromFile(cert, "")
	if err != nil {
		t.Fatalf("Failed to load certificate file: %v", err)
	}
	conn, err := grpc.Dial(adress,
		grpc.WithTransportCredentials(tlsCreds),
		tutorial.WithClientInterceptor(),
	)
	if err != nil {
		t.Fatalf("Failed to dial server: %v", err)
	}
	client := api.NewCacheClient(conn)

	// test
	fn(client)
}

func TestBadToken(t *testing.T) {
	withClient(t, func(client api.CacheClient) {
		// no token
		ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
		_, err := client.Store(ctx, &api.StoreReq{
			Key: "gopher",
			Val: []byte("con"),
		})
		if err == nil {
			t.Fatal("Store succeeded, expect failed")
		}
		if _, ok := status.FromError(err); !ok {
			t.Fatalf("Got unknown error: %v", err)
		}

		// bad token
		ctx, _ = context.WithTimeout(context.Background(), 100*time.Millisecond)
		_, err = client.Store(ctx, &api.StoreReq{
			AccountToken: "foo",
			Key:          "gopher",
			Val:          []byte("con"),
		})
		if err == nil {
			t.Fatal("Store succeeded, expect failed")
		}
		if _, ok := status.FromError(err); !ok {
			t.Fatalf("Got unknown error: %v", err)
		}
	})
}

func TestBasic(t *testing.T) {
	withClient(t, func(client api.CacheClient) {
		// store
		ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
		_, err := client.Store(ctx, &api.StoreReq{
			AccountToken: "token",
			Key:          "gopher",
			Val:          []byte("con"),
		})
		if err != nil {
			t.Fatalf("Failed to store: %v", err)
		}

		// get
		ctx, _ = context.WithTimeout(context.Background(), 10*time.Millisecond)
		resp, err := client.Get(ctx, &api.GetReq{Key: "gopher"})
		if err != nil {
			t.Fatalf("Failed to get: %v", err)
		}
		if string(resp.Val) != "con" {
			t.Fatalf("Got cached value: %s, want: %s", resp.Val, "con")
		}
	})
}

func TestDryRun(t *testing.T) {
	withClient(t, func(client api.CacheClient) {
		// store (dry run)
		ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("dry-run", "1"))
		_, err := client.Store(ctx, &api.StoreReq{
			AccountToken: "token",
			Key:          "con",
			Val:          []byte("2017"),
		})
		if err != nil {
			t.Fatalf("Failed to store: %v", err)
		}

		// get, expects not found
		ctx, _ = context.WithTimeout(context.Background(), 10*time.Millisecond)
		resp, err := client.Get(ctx, &api.GetReq{Key: "con"})
		if err == nil {
			t.Fatalf("Got cached value: %s", resp.Val)
		}
		if _, ok := status.FromError(err); !ok {
			t.Fatalf("Got unknown error: %v", err)
		}
	})
}

func TestDump(t *testing.T) {
	withClient(t, func(client api.CacheClient) {
		// store 1
		ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
		_, err := client.Store(ctx, &api.StoreReq{
			AccountToken: "token",
			Key:          "gopher",
			Val:          []byte("con"),
		})
		if err != nil {
			t.Fatalf("Failed to store: %v", err)
		}

		// store 2
		ctx, _ = context.WithTimeout(context.Background(), 100*time.Millisecond)
		_, err = client.Store(ctx, &api.StoreReq{
			AccountToken: "token",
			Key:          "go",
			Val:          []byte("func"),
		})
		if err != nil {
			t.Fatalf("Failed to store: %v", err)
		}

		// dump
		ctx, _ = context.WithTimeout(context.Background(), 10*time.Millisecond)
		dump, err := client.Dump(ctx, &api.DumpReq{})
		if err != nil {
			t.Fatalf("Failed to dump: %v", err)
		}
		if len(dump.Items) != 2 {
			t.Fatalf("Items len: %d, want: 2", len(dump.Items))
		}
	})
}
