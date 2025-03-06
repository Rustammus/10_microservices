package grps

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"microservices/task_5/server/internal/models"
	"microservices/task_5/server/proto/auth"
	"strings"
)

type claimsKey string

const ClaimsKey claimsKey = "claims"

var (
	errMissingMetadata  = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken     = status.Errorf(codes.Unauthenticated, "invalid token")
	errPermissionDenied = status.Errorf(codes.PermissionDenied, "permission denied")
)

func (s *RPCServer) AuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	// method not require authorization
	if info.FullMethod == auth.Auth_Register_FullMethodName || info.FullMethod == auth.Auth_Login_FullMethodName {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	claims, err := s.parse(ctx, md["authorization"])
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, ClaimsKey, claims)

	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}

// valid validates the authorization.
func (s *RPCServer) parse(_ context.Context, authorization []string) (models.Claims, error) {
	if len(authorization) < 1 {
		return models.Claims{}, errMissingMetadata
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")

	claims, err := s.s.Auth.Authenticate(nil, models.Token{AccessToken: token})

	if err != nil {
		return models.Claims{}, err
	}

	return claims, nil
}
