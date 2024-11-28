package subscriber

import (
	"context"
	"log"
	"sync"

	"cloud.google.com/go/pubsub"
)

type Subscriber struct {
	client *pubsub.Client
}

var once sync.Once
var sub Subscriber

func GetClient() Subscriber {
	return sub
}

func InitClient(ctx context.Context, projectId string) {
	once.Do(func() {
		client, err := pubsub.NewClient(ctx, projectId)
		if err != nil {
			panic(err)
		}
		sub = Subscriber{client: client}
	})
}

// Subscribe subscribes to a topic and uses the handler func arg to
// process the message.
func (s *Subscriber) Subscribe(ctx context.Context, subId string, handler func(context.Context, *pubsub.Message) error) {
	subscription := s.client.Subscription(subId)
	err := subscription.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		err := handler(ctx, m)
		if err != nil {
			log.Printf("Failed to proccess message. %v", err)
		} else {
			m.Ack()
		}
	})
	if err != nil {
		log.Printf("Failed to receive message. %v\n", err)
	}
}
