package app

import (
	"context"
	"encoding/json"
	ak "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
	"math/rand"
	"microservices/task_5/emailServiceMock/internal/schemas"
	"microservices/task_5/emailServiceMock/pkg/kafka"
	"os/signal"
	"syscall"
	"time"

	"os"
)

var confirmTopic = "notify-confirm"

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

	p := kafka.NewProducer(bHost + ":" + bPort)

	c, err := kafka.NewConsumer("notify-topic", bHost+":"+bPort, "myGroup")
	if err != nil {
		panic(err)
	}

	go emailSendEmulation(p, c)

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()
}

func emailSendEmulation(p *kafka.Producer, c *kafka.Consumer) {

	for msg := range c.MsgChan() {
		notify := &schemas.NotifySend{}
		err := json.Unmarshal(msg.Value, notify)
		if err != nil {
			log.Println(err)
		}
		log.Printf("receive notify %d", notify.Id)

		go func(id int64) {
			confirm := &schemas.NotifyConfirm{Id: id}
			msgRaw, err := json.Marshal(confirm)
			if err != nil {
				log.Println(err)
			}

			// work...
			tts := time.Duration(rand.Intn(20)+10) * time.Second
			time.Sleep(tts)

			p.Send(&ak.Message{
				TopicPartition: ak.TopicPartition{
					Topic:     &confirmTopic,
					Partition: ak.PartitionAny,
				},
				Value: msgRaw,
			})
			log.Printf("send confirm %d", notify.Id)
		}(notify.Id)
	}
}
