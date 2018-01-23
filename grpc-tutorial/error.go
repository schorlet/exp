package tutorial

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (e *Error) Error() string {
	return e.Message
}

// Errorf returns *rpc.Error.
func Errorf(code codes.Code, temporary bool, msg string, args ...interface{}) error {
	return &Error{
		Code:      int64(code),
		Message:   fmt.Sprintf(msg, args...),
		Temporary: temporary,
	}
}

// MarshalError sets err in the trailer metadata.
// Returns a RPC status code (status.*Status).
// Does nothing if err can't be casted to *rpc.Error.
func MarshalError(err error, ctx context.Context) error {
	rerr, ok := err.(*Error)
	if !ok {
		return err
	}

	pberr, marshalerr := proto.Marshal(rerr)
	if marshalerr == nil {
		md := metadata.Pairs("rpc-error", base64.StdEncoding.EncodeToString(pberr))
		trailerr := grpc.SetTrailer(ctx, md)
		if trailerr != nil {
			log.Printf("Unable to set rpc-error metadata: %v\n", trailerr)
		}
	}

	return status.Errorf(codes.Code(rerr.Code), rerr.Message)
}

// UnmarshalError gets an *rpc.Error from metadata, if any.
func UnmarshalError(md metadata.MD) *Error {
	vals, ok := md["rpc-error"]
	if !ok {
		return nil
	}
	if len(vals) < 1 {
		return nil
	}
	buf, err := base64.StdEncoding.DecodeString(vals[0])
	if err != nil {
		return nil
	}
	var rerr Error
	if err := proto.Unmarshal(buf, &rerr); err != nil {
		return nil
	}
	return &rerr
}
