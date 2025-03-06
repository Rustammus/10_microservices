package app

import (
	"awesomeProject1/internal/controllers/grps"
	"awesomeProject1/internal/controllers/rest"
	"awesomeProject1/internal/crud"
	"awesomeProject1/internal/services"
	"awesomeProject1/internal/services/auth"
	"awesomeProject1/internal/services/user"
	"log"
	"os"
	"sync"
	"time"
)

func Run() {

	userCrud := crud.NewUserCRUD()

	userService := user.NewService(userCrud)

	authService := auth.NewService(os.Getenv("APP_SALT"), time.Hour*24)

	s := &services.Services{
		User: userService,
		Auth: authService,
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	api := rest.NewAPI(s)
	go api.Run(wg)

	wg.Add(1)
	rpc := grps.NewRPCServer(s)
	go rpc.Run(wg)

	wg.Wait()
	log.Println("stop")
}
