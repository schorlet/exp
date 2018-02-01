package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"

	tutorial "github.com/schorlet/exp/grpc-tutorial"

	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
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
	srv := grpc.NewServer(
		grpc.Creds(tlsCreds),
		tutorial.ServerInterceptor(),
	)

	cacheService := CacheService{
		store: map[string][]byte{},
		// accounts:,
		keysByAccount: map[string]int64{},
	}
	tutorial.RegisterCacheServer(srv, &cacheService)

	tutorial.RegisterAccountsServer(srv, &AccountsService{
		store: map[string]tutorial.Account{
			"token": {MaxCacheKeys: 2},
		},
	})

	l, err := net.Listen("tcp", "localhost:5051")
	if err != nil {
		return err
	}

	var g errgroup.Group
	g.Go(func() error {
		// blocks until complete
		return srv.Serve(l)
	})
	g.Go(func() error {
		tlsClient := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
		conn, err := grpc.Dial("localhost:5051",
			grpc.WithTransportCredentials(tlsClient),
			tutorial.WithClientInterceptor(),
		)
		if err != nil {
			return fmt.Errorf("failed to dial server: %v", err)
		}
		cacheService.accounts = tutorial.NewAccountsClient(conn)
		return nil
	})

	return g.Wait()
}

// CacheService stores values in memory.
type CacheService struct {
	store         map[string][]byte
	accounts      tutorial.AccountsClient
	keysByAccount map[string]int64
}

// Get returns a value from the cache.
func (s *CacheService) Get(ctx context.Context, req *tutorial.GetReq) (*tutorial.GetResp, error) {
	val, ok := s.store[req.Key]
	if !ok {
		return nil, tutorial.Errorf(codes.NotFound, true, "key not found %q", req.Key)
	}
	return &tutorial.GetResp{Val: val}, nil
}

// Store sets a value into the cache.
func (s *CacheService) Store(ctx context.Context, req *tutorial.StoreReq) (*tutorial.StoreResp, error) {
	// ctx is propagated from the original client call through all sub services calls,
	// so is the deadline, the timeout for how long the entire operation takes
	resp, err := s.accounts.GetByToken(ctx,
		&tutorial.GetByTokenReq{Token: req.AccountToken})
	if err != nil {
		return nil, status.Errorf(codes.Unknown,
			"failed to get token %q: %v", req.AccountToken, err)
	}

	if s.keysByAccount[req.AccountToken] >= resp.Account.MaxCacheKeys {
		return nil, status.Errorf(codes.FailedPrecondition,
			"account %q exceeds max key limit %d", req.AccountToken, resp.Account.MaxCacheKeys)
	}

	if !dryRun(ctx) {
		if _, ok := s.store[req.Key]; !ok {
			s.keysByAccount[req.AccountToken]++
		}
		s.store[req.Key] = req.Val
	}

	return &tutorial.StoreResp{}, nil
}

func dryRun(ctx context.Context) bool {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false
	}

	val, ok := md["dry-run"]
	if !ok {
		return false
	}

	if len(val) < 1 {
		return false
	}

	return val[0] == "1"
}

// Dump returns all values from the cache.
func (s *CacheService) Dump(context.Context, *tutorial.DumpReq) (*tutorial.DumpResp, error) {
	var i int
	resp := tutorial.DumpResp{Items: make([]*tutorial.DumpItem, len(s.store))}
	for key, val := range s.store {
		resp.Items[i] = &tutorial.DumpItem{Key: key, Val: val}
		i++
	}
	return &resp, nil
}

// AccountsService stores Accounts in memory.
type AccountsService struct {
	store map[string]tutorial.Account
}

// GetByToken returns an Account.
func (a *AccountsService) GetByToken(ctx context.Context, req *tutorial.GetByTokenReq) (*tutorial.GetByTokenResp, error) {
	val, ok := a.store[req.Token]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "token not found %q", req.Token)
	}
	return &tutorial.GetByTokenResp{Account: &val}, nil
}
