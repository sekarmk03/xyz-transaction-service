package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	roles "xyz-transaction-service/common/authorization"
	commonJwt "xyz-transaction-service/common/jwt"
	"xyz-transaction-service/server/interceptor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	connProtocol  = "tcp"
	maxMsgSize    = 1024 * 1024 * 150
	tokenDuration = 5 * time.Minute
	secretKey     = "secret"
)

type Grpc struct {
	Server   *grpc.Server
	listener net.Listener
	Port     string
}

func NewGrpc(port string, options ...grpc.ServerOption) *Grpc {
	options = append(options, grpc.MaxSendMsgSize(maxMsgSize))
	options = append(options, grpc.MaxRecvMsgSize(maxMsgSize))

	server := grpc.NewServer(options...)

	return &Grpc{
		Server: server,
		Port:   port,
	}
}

func NewGrpcServer(port string, jwtManager *commonJwt.JWT) *Grpc {
	authInterceptor := interceptor.NewAuthInterceptor(jwtManager, roles.GetAccessibleRoles())
	options := []grpc.ServerOption{
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	}

	server := NewGrpc(port, options...)
	return server
}

func (g *Grpc) Run() error {
	var err error
	g.listener, err = net.Listen(connProtocol, fmt.Sprintf(":%s", g.Port))
	if err != nil {
		return status.Errorf(codes.Internal, "ERROR: Failed to listen on port %s: %v", g.Port, err)
	}

	go g.serve()
	log.Printf("grpc server is running on port %s\n", g.Port)
	return nil
}

func (g *Grpc) serve() {
	if err := g.Server.Serve(g.listener); err != nil {
		panic(err)
	}
}

func (g *Grpc) AwaitTermination() error {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	<-sign

	g.Server.GracefulStop()
	return g.listener.Close()
}
