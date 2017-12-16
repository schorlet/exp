package interceptor

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

// ServerInterceptor returns a ServerOption which logs RPC calls on stderr.
func ServerInterceptor() grpc.ServerOption {
	// Only one unary interceptor can be installed.
	// The construction of multiple interceptors (e.g., chaining) can be implemented at the caller.
	return grpc.UnaryInterceptor(serverInterceptor)
}

func serverInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now()
	resp, err = handler(ctx, req)
	log.Printf("invoke server method=%q duration=%s error=%v",
		info.FullMethod, time.Since(start), err)
	return
}
