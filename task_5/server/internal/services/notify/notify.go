package notify

import (
	"encoding/json"
	ak "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
	"microservices/task_5/server/internal/models"
	"microservices/task_5/server/internal/schemas"
	"sync/atomic"
)

type Producer interface {
	Close()
	Send(msg *ak.Message)
}

type Consumer interface {
	Close()
	MsgChan() (m <-chan *ak.Message)
}

type Service struct {
	counter atomic.Int64
	topic   string
	p       Producer
	c       Consumer
}

func NewService(topic string, p Producer, c Consumer) *Service {
	s := &Service{topic: topic, p: p, c: c}
	go s.receiveConfirmation()
	return s
}

func (s *Service) UserCreated(user models.User) {
	msg := &schemas.NotifySend{
		UserCreated: &schemas.NotifyUserCreated{User: user},
	}

	s.send(msg)
}

func (s *Service) UserUpdated(userOld, userNew models.User) {
	msg := &schemas.NotifySend{
		UserUpdated: &schemas.NotifyUserUpdated{UserOld: userOld, UserNew: userNew},
	}

	s.send(msg)
}

func (s *Service) UserDeleted(user models.User) {
	msg := &schemas.NotifySend{
		UserDeleted: &schemas.NotifyUserDeleted{User: user},
	}

	s.send(msg)
}

func (s *Service) send(notify *schemas.NotifySend) {
	notify.Id = s.counter.Load()
	msgRaw, err := json.Marshal(notify)
	if err != nil {
		log.Print("notifyService: marshal err: ", err)
		return
	}

	s.p.Send(&ak.Message{
		TopicPartition: ak.TopicPartition{
			Topic:     &(s.topic),
			Partition: ak.PartitionAny,
		},
		Value: msgRaw,
	})

	log.Printf("notifyService: send msg: %+v", notify)
	s.counter.Add(1)
}

func (s *Service) receiveConfirmation() {

	for msg := range s.c.MsgChan() {
		confiration := &schemas.NotifyConfirm{}
		if err := json.Unmarshal(msg.Value, confiration); err != nil {
			log.Print("notifyService: unmarshal err: ", err)
		}
		log.Printf("notifyService: notify %d delivered", confiration.Id)
	}

}
