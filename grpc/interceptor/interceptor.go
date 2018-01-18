package interceptor

import (
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WithClientInterceptor returns a DialOption which logs RPC calls on stderr.
func WithClientInterceptor() grpc.DialOption {
	return grpc.WithUnaryInterceptor(clientInterceptor)
}

// NotIdempotent returns a copy of ctx with "idempotent" valued to false.
func NotIdempotent(ctx context.Context) context.Context {
	return context.WithValue(ctx, "idempotent", false)
}

// IsIdempotent returns the "idempotent" value stored in ctx, if any.
func IsIdempotent(ctx context.Context) bool {
	val, ok := ctx.Value("idempotent").(bool)
	if !ok {
		return true
	}
	return val
}

func clientInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	var (
		attempts int
		err      error
	)
	for attempts < 3 {
		select {
		case <-ctx.Done():
			err = status.Errorf(codes.DeadlineExceeded,
				"timeout reached after %d attempts: %v", attempts, ctx.Err())
		default:
			attempts++
			start := time.Now()
			err = invoker(ctx, method, req, reply, cc, opts...)
			log.Printf("invoke=%d remote method=%q duration=%s error=%v",
				attempts, method, time.Since(start), err)
			if err != nil && IsIdempotent(ctx) {
				continue
			}
		}
		break
	}
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
