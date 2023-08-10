package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wisesight/go-api-template/config"
	"github.com/wisesight/go-api-template/pkg/adapter"
	"github.com/wisesight/go-api-template/pkg/entity"
	"github.com/wisesight/go-api-template/pkg/validator"
	"go.mongodb.org/mongo-driver/bson"
)

type hashtag struct {
	mongoDBAdapter    adapter.IMongoDBAdapter
	hashtagCollection adapter.IMongoCollection
	timeout           time.Duration
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func (h hashtag) GetAll() ([]entity.Hashtag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var hashtags []entity.Hashtag

	err := h.mongoDBAdapter.Find(ctx, h.hashtagCollection, &hashtags, bson.D{})

	failOnError(err, "Failed to connect to MongoDB")

	return hashtags, nil
}

func main() {
	cfg := config.NewConfig()

	fmt.Println(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongodbClient, err := adapter.NewMongoDBConnection(ctx, "mongodb://root:root@localhost:27017")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = mongodbClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	if err = validator.NewValidator(); err != nil {
		panic(err)
	}

	hashtagCollection := mongodbClient.Database("go_workshop").Collection("hashtags")
	mongoDBAdapter := adapter.NewMongoDBAdapter(mongodbClient)

	h := hashtag{hashtagCollection: hashtagCollection, mongoDBAdapter: mongoDBAdapter}

	res, err := h.GetAll()

	if err != nil {
		fmt.Println("Error")
	}

	conn, err := amqp.Dial("amqp://root:root@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to connect to Channel")
	defer ch.Close()

	// declare queue
	q, err := ch.QueueDeclare("local:workshop:hashtag:job", false, false, false, false, nil)
	failOnError(err, "Failed to declare queue")

	for i := 0; i < len(res); i++ {
		body, _ := json.Marshal(res[i])
		err = ch.PublishWithContext(ctx,
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s\n", body)
	}

}
