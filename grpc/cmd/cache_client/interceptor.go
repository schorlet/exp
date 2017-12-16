package main

import (
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// WithClientInterceptor returns a DialOption which logs RPC calls on stderr.
func WithClientInterceptor() grpc.DialOption {
	return grpc.WithUnaryInterceptor(clientInterceptor)
}

func clientInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("invoke remote method=%q duration=%s error=%v",
		method, time.Since(start), err)
	return err
}
