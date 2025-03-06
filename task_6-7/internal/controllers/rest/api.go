package rest

import (
	"context"
	"errors"
	"github.com/julienschmidt/httprouter"
	"log/slog"
	"microservices/task_6/internal/config"
	"microservices/task_6/internal/metrics"
	"microservices/task_6/internal/services"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type API struct {
	s *services.Services
	l *slog.Logger
	m *metrics.AppMetrics
}

func NewAPI(s *services.Services, l *slog.Logger, m *metrics.AppMetrics) *API {
	l = l.With(slog.String("api", "rest"))
	return &API{s: s, l: l, m: m}
}

func (s *API) Run(wg *sync.WaitGroup) {
	router := httprouter.New()
	s.initUserController(router)

	c := config.GetConfig()

	server := http.Server{
		Addr:    ":" + c.REST.Port,
		Handler: s.BaseMiddleware(s.LoggerMiddleware(router)),
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	serverShutdown := sync.WaitGroup{}
	go func() {
		sig := <-sigChan

		serverShutdown.Add(1)

		s.l.Info("got interrupt signal, shutting down server", slog.String("os_signal", sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			s.l.Error("error on shutting down server", slog.String("error", err.Error()))
		} else {
			s.l.Info("server shutdown gracefully")
		}

		serverShutdown.Done()
	}()

	// run server
	s.l.Info("try to run server")
	err := server.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			s.l.Error("listen error", slog.String("error", err.Error()))
		}
	}

	// Wait for server Shutdown
	serverShutdown.Wait()

	signal.Stop(sigChan)
	close(sigChan)

	wg.Done()
}
