package server

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DialOption func(name string) (grpc.DialOption, error)

func Dial(name string, opts ...DialOption) (*grpc.ClientConn, error) {
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	for _, fn := range opts {
		opt, err := fn(name)
		if err != nil {
			return nil, fmt.Errorf("config error: %v", err)
		}
		dialOpts = append(dialOpts, opt)
	}

	conn, err := grpc.Dial(name, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %v", name, err)
	}

	return conn, nil
}

func InitGRPCConn(addr string, ssl bool, cert string) *grpc.ClientConn {
	// if ssl true, dial with ssl

	// else
	conn, err := Dial(addr)
	if err != nil {
		panic(fmt.Sprintf("ERROR: dial error: %v", err))
	}

	return conn
}