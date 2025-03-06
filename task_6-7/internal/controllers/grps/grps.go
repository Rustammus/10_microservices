package grps

import (
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"microservices/task_6/internal/config"
	"microservices/task_6/internal/metrics"
	"microservices/task_6/internal/services"
	"microservices/task_6/proto/auth"
	"microservices/task_6/proto/user"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type RPCServer struct {
	s *services.Services
	l *slog.Logger
	m *metrics.AppMetrics
}

func NewRPCServer(s *services.Services, l *slog.Logger, m *metrics.AppMetrics) *RPCServer {
	l = l.With(slog.String("api", "grps"))
	return &RPCServer{s: s, l: l, m: m}
}

func (s *RPCServer) Run(wg *sync.WaitGroup) {
	c := config.GetConfig()

	lis, err := net.Listen("tcp", ":"+c.RPC.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(s.LoggerMiddleware(s.AuthInterceptor)))

	user.RegisterUserServer(grpcServer, &userService{s: s.s})
	auth.RegisterAuthServer(grpcServer, &authService{s: s.s})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		sig := <-sigChan
		s.l.Info("got interrupt signal, shutting down server", slog.String("os_signal", sig.String()))

		grpcServer.GracefulStop()

		s.l.Info("server shutdown gracefully")
	}()

	s.l.Info("try to run server")
	if err := grpcServer.Serve(lis); err != nil {
		s.l.Error("listen error", slog.String("error", err.Error()))
	}

	signal.Stop(sigChan)
	close(sigChan)

	wg.Done()
}
