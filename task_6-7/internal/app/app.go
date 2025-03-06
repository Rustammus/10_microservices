package app

import (
	"log"
	"log/slog"
	"microservices/task_6/internal/config"
	"microservices/task_6/internal/controllers/grps"
	"microservices/task_6/internal/controllers/rest"
	"microservices/task_6/internal/crud"
	"microservices/task_6/internal/metrics"
	"microservices/task_6/internal/models"
	"microservices/task_6/internal/services"
	"microservices/task_6/internal/services/auth"
	"microservices/task_6/internal/services/user"
	"microservices/task_6/pkg/logger"
	"os"
	"runtime"
	"sync"
	"time"
)

func Run() {
	// Config init
	c := config.GetConfig()

	// Logger init
	l := logger.NewLogger(c.Log.Level, c.Log.Output)
	logOsInfo(l)

	u := &models.User{
		ID:       777777,
		Name:     "Telegram",
		Email:    "tg@gmail.com",
		Password: "404",
	}

	l.Info("test", slog.Any("user", u))

	// Metrics service init
	m := metrics.NewAppMetrics()
	defer m.Close()

	// Repository init
	userCrud := crud.NewUserCRUD(l)

	// Services init
	userService := user.NewService(userCrud)

	authService := auth.NewService("salt", time.Hour*24)

	s := &services.Services{
		User: userService,
		Auth: authService,
	}

	wg := &sync.WaitGroup{}

	// REST server init
	wg.Add(1)
	api := rest.NewAPI(s, l, m)
	go api.Run(wg)

	// gRPS server init
	wg.Add(1)
	rpc := grps.NewRPCServer(s, l, m)
	go rpc.Run(wg)

	wg.Wait()
	log.Println("stop")
}

func logOsInfo(l *slog.Logger) {
	workdir, _ := os.Getwd()

	gomaxprocs := runtime.GOMAXPROCS(0)

	l.Info("server info", slog.String("workdir", workdir), slog.Int("go_max_procs", gomaxprocs))
}
