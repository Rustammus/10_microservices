package grps

import (
	"google.golang.org/grpc"
	"log"
	"microservices/task_5/server/internal/services"
	"microservices/task_5/server/proto/auth"
	"microservices/task_5/server/proto/user"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type RPCServer struct {
	s *services.Services
}

func NewRPCServer(s *services.Services) *RPCServer {
	return &RPCServer{s: s}
}

func (s *RPCServer) Run(wg *sync.WaitGroup) {

	port, ok := os.LookupEnv("APP_RPC_PORT")
	if !ok {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(s.AuthInterceptor))
	//grpcServer := grpc.NewServer()
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
		log.Print("gRPC: Shutdown signal: ", sig)
		grpcServer.GracefulStop()
		log.Println("gRPC: Gracefully stopped")
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("failed to serve: %v\n", err)
	}

	signal.Stop(sigChan)
	close(sigChan)

	wg.Done()
}
