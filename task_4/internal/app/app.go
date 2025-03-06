package app

import (
	"log"
	"microservices/task_4/internal/controllers/rest"
	"microservices/task_4/internal/services"
	"microservices/task_4/internal/services/chat"
	"sync"
)

func Run() {

	s := &services.Services{
		Chat: chat.NewService(),
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	api := rest.NewAPI(s)
	go api.Run(wg)

	wg.Wait()
	log.Println("stop")
}
