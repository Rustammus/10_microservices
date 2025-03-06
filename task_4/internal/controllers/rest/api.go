package rest

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"log"
	"microservices/task_4/internal/services"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type API struct {
	s *services.Services
}

func NewAPI(s *services.Services) *API {
	return &API{s: s}
}

// Run server
func (s *API) Run(wg *sync.WaitGroup) {

	// init routes
	router := httprouter.New()
	s.initChatController(router)

	port, ok := os.LookupEnv("APP_REST_PORT")
	if !ok {
		port = "8080"
	}

	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
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
		log.Print("rest: Shutdown signal: ", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			log.Print("rest: Error shutting down: ", err)
		} else {
			log.Println("rest: Gracefully stopped")
		}
		serverShutdown.Done()
	}()

	log.Print(server.ListenAndServe())

	serverShutdown.Wait()

	signal.Stop(sigChan)
	close(sigChan)

	wg.Done()
}
