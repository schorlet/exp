package tutorial

import (
	"fmt"
	"net"

	"github.com/schorlet/exp/grpc-tutorial/api"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Serve starts a new gRPC server.
func Serve(address, certFile, keyFile string) (func(), error) {
	tlsServer, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate file: %v", err)
	}
	server := grpc.NewServer(
		grpc.Creds(tlsServer),
		ServerInterceptor(),
	)

	cacheService := CacheService{
		store: map[string][]byte{},
		// accounts: is api.AccountsClient,
		keysByAccount: map[string]int64{},
	}
	accountsService := AccountsService{
		store: map[string]api.Account{
			"token": {MaxCacheKeys: 2},
		},
	}
	api.RegisterCacheServer(server, &cacheService)
	api.RegisterAccountsServer(server, &accountsService)

	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err)
	}
	// start server, blocks until complete
	go server.Serve(l)

	// create accounts client
	tlsClient, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate file: %v", err)
	}
	conn, err := grpc.Dial(address,
		grpc.WithTransportCredentials(tlsClient),
		WithClientInterceptor(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial server: %v", err)
	}
	cacheService.accounts = api.NewAccountsClient(conn)

	return server.Stop, nil
}

// CacheService stores values in memory.
type CacheService struct {
	store         map[string][]byte
	accounts      api.AccountsClient
	keysByAccount map[string]int64
}

// Get returns a value from the cache.
func (s *CacheService) Get(ctx context.Context, req *api.GetReq) (*api.GetResp, error) {
	val, ok := s.store[req.Key]
	if !ok {
		return nil, api.Errorf(codes.NotFound, true, "key not found %q", req.Key)
	}
	return &api.GetResp{Val: val}, nil
}

// Store sets a value into the cache.
func (s *CacheService) Store(ctx context.Context, req *api.StoreReq) (*api.StoreResp, error) {
	// ctx is propagated from the original client call through all sub services calls,
	// so is the deadline, the timeout for how long the entire operation takes
	resp, err := s.accounts.GetByToken(ctx,
		&api.GetByTokenReq{Token: req.AccountToken})
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

	return &api.StoreResp{}, nil
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
func (s *CacheService) Dump(context.Context, *api.DumpReq) (*api.DumpResp, error) {
	var i int
	resp := api.DumpResp{Items: make([]*api.DumpItem, len(s.store))}
	for key, val := range s.store {
		resp.Items[i] = &api.DumpItem{Key: key, Val: val}
		i++
	}
	return &resp, nil
}

// AccountsService stores Accounts in memory.
type AccountsService struct {
	store map[string]api.Account
}

// GetByToken returns an Account.
func (a *AccountsService) GetByToken(ctx context.Context, req *api.GetByTokenReq) (*api.GetByTokenResp, error) {
	val, ok := a.store[req.Token]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "token not found %q", req.Token)
	}
	return &api.GetByTokenResp{Account: &val}, nil
}
