package app

import (
	"log"
	"microservices/task_5/server/internal/controllers/grps"
	"microservices/task_5/server/internal/controllers/rest"
	"microservices/task_5/server/internal/crud"
	"microservices/task_5/server/internal/services"
	"microservices/task_5/server/internal/services/auth"
	"microservices/task_5/server/internal/services/notify"
	"microservices/task_5/server/internal/services/user"
	"microservices/task_5/server/pkg/kafka"
	"os"
	"sync"
	"time"
)

func Run() {
	// read envs
	bHost, ok := os.LookupEnv("KAFKA_HOST")
	if !ok {
		bHost = "localhost"
	}
	log.Println("KAFKA_HOST " + bHost)

	bPort, ok := os.LookupEnv("KAFKA_PORT")
	if !ok {
		bPort = "9092"
	}
	log.Println("KAFKA_PORT " + bPort)

	// Create producer (send notify to email service)
	p := kafka.NewProducer(bHost + ":" + bPort)
	defer p.Close()

	// Create consumer (receive confirmation from email service)
	c, err := kafka.NewConsumer("notify-confirm", bHost+":"+bPort, "myGroup")
	if err != nil {
		log.Printf("Failed to create consumer: %v", err)
		return
	}
	defer c.Close()

	notifyService := notify.NewService("notify-topic", p, c)

	userCrud := crud.NewUserCRUD()

	userService := user.NewService(userCrud, notifyService)

	authService := auth.NewService("salt", time.Hour*24)

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
