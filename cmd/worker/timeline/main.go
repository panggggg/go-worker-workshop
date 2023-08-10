package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wisesight/go-api-template/pkg/entity"
	"github.com/wisesight/go-api-template/pkg/repository"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
		return
	}
}

func main() {
	conn, err := amqp.Dial("amqp://root:root@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to connect to Channel")
	defer ch.Close()

	// declare queue
	q, err := ch.QueueDeclare("local:workshop:hashtag:result", false, false, false, false, nil)
	failOnError(err, "Failed to declare queue")

	jobs, err := ch.Consume("local:workshop:hashtag:job", "", true, false, false, false, nil)
	failOnError(err, "Failed to consume jobs")

	var forever chan struct{}

	// marshal: struct -> json
	// unmarshall: json -> struct

	var job entity.Hashtag
	var messages []*entity.Job
	var mutex sync.Mutex
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for j := range jobs {
		log.Printf("Received a job: %s \n", j.Body)
		json.Unmarshal(j.Body, &job)
		threads, _ := repository.GetThreads(job.Keyword)
		for _, data := range threads.Data {
			wg.Add(1)
			go func(data entity.Thread) {
				defer wg.Done()
				result, err := repository.GetAccountInfo(data)
				failOnError(err, "Cannot get account info")
				mutex.Lock()
				messages = append(messages, result)
				fmt.Println("MESSAGE ===> ", messages)
				mutex.Unlock()
			}(data)
			wg.Wait()
		}
		for _, m := range messages {
			body, _ := json.Marshal(*m)
			fmt.Println("BODY =====> ", string(body))
			err = ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte(body)})
			fmt.Println("Send message to queue")
		}
	}

	log.Printf(" [*] Waiting for jobs. To exit press CTRL+C")
	<-forever
}
