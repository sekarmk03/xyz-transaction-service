package utils

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func GetMetadataAuthorization(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("ERROR: [Utils - GetMetadataAuthorization] Metadata is not provided")
		return "", status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values, ok := md["authorization"]
	if !ok || len(values) == 0 {
		log.Println("ERROR: [Utils - GetMetadataAuthorization] Authorization token is not provided")
		return "", status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	authHeader := values[0]

	return authHeader, nil
}
