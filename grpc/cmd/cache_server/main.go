package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/schorlet/exp/grpc/rpc"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
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
	srv := grpc.NewServer(grpc.Creds(tlsCreds))

	cacheService := CacheService{
		store: map[string][]byte{},
		// accounts:,
		keysByAccount: map[string]int64{},
	}
	rpc.RegisterCacheServer(srv, &cacheService)

	rpc.RegisterAccountsServer(srv, &AccountsService{
		store: map[string]rpc.Account{
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
		conn, err := grpc.Dial("localhost:5051", grpc.WithTransportCredentials(tlsClient))
		if err != nil {
			return fmt.Errorf("failed to dial server: %v", err)
		}
		cacheService.accounts = rpc.NewAccountsClient(conn)
		return nil
	})

	return g.Wait()
}

// CacheService stores values in memory.
type CacheService struct {
	store         map[string][]byte
	accounts      rpc.AccountsClient
	keysByAccount map[string]int64
}

// Get returns a value from the cache
func (s *CacheService) Get(ctx context.Context, req *rpc.GetReq) (*rpc.GetResp, error) {
	val, ok := s.store[req.Key]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Key not found %q", req.Key)
	}
	return &rpc.GetResp{Val: val}, nil
}

// Store sets a value into the cache
func (s *CacheService) Store(ctx context.Context, req *rpc.StoreReq) (*rpc.StoreResp, error) {
	start := time.Now()
	resp, err := s.accounts.GetByToken(context.Background(),
		&rpc.GetByTokenReq{Token: req.AccountToken})
	log.Printf("accounts.GetByToken duration %s", time.Since(start))
	if err != nil {
		return nil, status.Errorf(codes.Unknown,
			"Failed to get token %q: %v", req.AccountToken, err)
	}

	if s.keysByAccount[req.AccountToken] >= resp.Account.MaxCacheKeys {
		return nil, status.Errorf(codes.FailedPrecondition,
			"Account %q exceeds max key limit %d", req.AccountToken, resp.Account.MaxCacheKeys)
	}

	if _, ok := s.store[req.Key]; !ok {
		s.keysByAccount[req.AccountToken]++
	}

	s.store[req.Key] = req.Val
	return &rpc.StoreResp{}, nil
}

// AccountsService stores Accounts in memory.
type AccountsService struct {
	store map[string]rpc.Account
}

// GetByToken returns an Account
func (a *AccountsService) GetByToken(ctx context.Context, req *rpc.GetByTokenReq) (*rpc.GetByTokenResp, error) {
	val, ok := a.store[req.Token]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Token not found %q", req.Token)
	}
	return &rpc.GetByTokenResp{Account: &val}, nil
}
