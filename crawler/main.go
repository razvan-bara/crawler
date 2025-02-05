package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"log"
	crawlerDb "main/db"
	"main/queue"
	"net/http"
	"os"
)

func main() {

	dbHost := getEnvOrPanic("DB_HOST")
	dbPort := getEnvOrPanic("DB_PORT")
	dbUser := getEnvOrPanic("DB_USER")
	dbPass := getEnvOrPanic("DB_PASSWORD")
	dbName := getEnvOrPanic("DB_NAME")

	db, err := crawlerDb.NewDb(dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}

	err = crawlerDb.DeleteAllArticles(db)
	if err != nil {
		log.Fatalf("failed to delete all articles: %v", err)
	}
	log.Println("Deleted all articles from the database")

	ampqHost := getEnvOrPanic("AMPQ_HOST")
	ampqPort := getEnvOrPanic("AMPQ_PORT")
	ampqVHost := getEnvOrPanic("AMPQ_VHOST")
	ampqUser := getEnvOrPanic("AMPQ_USER")
	ampqPass := getEnvOrPanic("AMPQ_PASS")

	connection, err := queue.ConnectToRabbitMQ(ampqUser, ampqPass, ampqHost, ampqPort, ampqVHost)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}

	messageQueue, err := queue.NewMessageQueue(connection, "crawler_queue")
	if err != nil {
		log.Fatalf("could not connect to queue: %v", err)
	}

	crawler, err := NewCrawler(
		messageQueue,
		db,
		&colly.LimitRule{
			DomainGlob:  "*",
			Parallelism: 1,
		},
		colly.AllowURLRevisit(),
		colly.MaxDepth(maxDepth),
		colly.UserAgent(userAgent),
	)
	if err != nil {
		log.Fatalf("could not create crawler: %v", err)
	}
	crawler.ConsumeTasks(messageQueue)

	isTstEnv := getEnvOrPanic("IS_TST")

	targetHost := targetRealUrl
	indexCount := dblpIndexCount
	if isTstEnv == "1" {
		targetHost = targetTestUrl
		indexCount = dblpTestIndexCount
	}

	r := mux.NewRouter()
	r.HandleFunc("/startCrawl", StartCrawl(messageQueue, targetHost, indexCount)).
		Methods("POST")

	initCrawlTask(messageQueue, targetHost, indexCount)

	addr := ":8081"
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func StartCrawl(messageQueue *queue.MessageQueue, targetHost string, indexCount int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		initCrawlTask(messageQueue, targetHost, indexCount)
		w.WriteHeader(http.StatusOK)
	}
}

func initCrawlTask(messageQueue *queue.MessageQueue, targetHost string, dblpIndexCount int) {
	for i := range dblpIndexCount {
		path := dblpIndex
		if i > 0 {
			path = fmt.Sprintf("%s?pos=%v", dblpIndex, 300*i+1)
		}
		url := targetHost + path
		log.Printf("Enqueued message for: %v", url)

		task := &Task{
			Url: url,
		}

		taskJson, err := json.Marshal(task)
		if err != nil {
			log.Fatalf("Failed to encode JSON message: %s", err)
		}

		err = messageQueue.Publish(amqp.Publishing{
			ContentType: "application/json",
			Body:        taskJson,
		})
		if err != nil {
			if errors.Is(err, queue.ErrorQueueMessageDuplicate) {
				log.Println("Not allowed to enqueue duplicate url message")
				continue
			}
			log.Printf("Failed to publish message: %v", err)
		}
	}
}

func getEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Environment variable %s not set", key))
	}

	return value
}
