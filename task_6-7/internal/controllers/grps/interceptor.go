package grps

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log/slog"
	"microservices/task_6/internal/models"
	"microservices/task_6/pkg/logger"
	"microservices/task_6/proto/auth"
	"strings"
	"time"
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

	ctx = logger.AppendCtx(ctx, slog.Int64("user_id", claims.UserID))

	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}

func (s *RPCServer) LoggerMiddleware(next grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

		// gen req_id into c
		reqId, err := uuid.NewV7()
		if err != nil {
			s.l.Error("error generate uuidV7", slog.String("error", err.Error()))
		} else {
			ctx = logger.AppendCtx(ctx, slog.String("req_id", reqId.String()))
		}

		// next handler
		tStart := time.Now()
		resp, err := next(ctx, req, info, handler)
		reqTime := time.Since(tStart)

		// error parse
		var errStr string
		if err != nil {
			errStr = err.Error()
		}

		// logging
		s.l.DebugContext(ctx, "new request",
			slog.String("method", info.FullMethod),
			slog.String("error", errStr),
			slog.String("duration", reqTime.String()))

		// Increment requests counter
		s.m.RequestsIncrement()

		return resp, err
	}
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
