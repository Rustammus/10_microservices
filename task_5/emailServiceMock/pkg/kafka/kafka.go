package kafka

import (
	"fmt"
	ak "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
	"runtime"
	"time"
)

type Consumer struct {
	msgChan   <-chan *ak.Message
	closeChan chan<- struct{}
}

func (c *Consumer) Close() {
	c.closeChan <- struct{}{}
	close(c.closeChan)
}

func (c *Consumer) MsgChan() (m <-chan *ak.Message) {
	return c.msgChan
}

func NewConsumer(topic, host, group string) (*Consumer, error) {

	_ = ak.NewTopicCollectionOfTopicNames([]string{topic})

	c, err := ak.NewConsumer(&ak.ConfigMap{
		"bootstrap.servers":        host,
		"group.id":                 group,
		"auto.offset.reset":        "earliest",
		"allow.auto.create.topics": true,
	})
	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		return nil, err
	}

	msgChan := make(chan *ak.Message)
	closeChan := make(chan struct{})
	go runConsume(c, msgChan, closeChan)

	return &Consumer{msgChan: msgChan, closeChan: closeChan}, nil
}

func runConsume(c *ak.Consumer, msgChan chan<- *ak.Message, closeChan <-chan struct{}) {
	for {
		select {

		case <-closeChan:
			err := c.Close()
			if err != nil {
				log.Println(err)
			}
			close(msgChan)

		default:
			msg, err := c.ReadMessage(time.Second)
			if err == nil {
				fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
				msgChan <- msg

			} else if !err.(ak.Error).IsTimeout() {
				// The client will automatically try to recover from all errors.
				// Timeout is not considered an error because it is raised by
				// ReadMessage in absence of messages.
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}

		// yield
		runtime.Gosched()
	}
}

type Producer struct {
	magChan chan<- *ak.Message
}

func (p *Producer) Close() {
	close(p.magChan)
}

func (p *Producer) Send(msg *ak.Message) {
	p.magChan <- msg
}

func NewProducer(host string) *Producer {
	p, err := ak.NewProducer(&ak.ConfigMap{"bootstrap.servers": host})
	if err != nil {
		panic(err)
	}

	msgChan := make(chan *ak.Message)

	go runProduce(p, msgChan)

	return &Producer{magChan: msgChan}
}

func runProduce(p *ak.Producer, msgChan <-chan *ak.Message) {
	events := p.Events()

forMark:
	for {
		select {
		case event := <-events:
			switch ev := event.(type) {
			case *ak.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
					ev.String()
				}
			}
		case msg, ok := <-msgChan:
			if !ok {
				// Wait for message deliveries before shutting down
				p.Flush(15 * 1000)

				p.Close()
				break forMark
			}

			err := p.Produce(msg, nil)
			if err != nil {
				fmt.Printf("Producer error: %v\n", err)
			}

		}
	}
}
