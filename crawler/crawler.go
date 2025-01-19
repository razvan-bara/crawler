package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/streadway/amqp"
	"log"
	crawlerDb "main/db"
	"main/queue"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	sleepTimeoutDuration   = time.Second * 5
	poolingTimeoutDuration = time.Second * 5
	limitHyperlinkPerPage  = 5
)

type Crawler struct {
	collector *colly.Collector
	queue     *queue.MessageQueue
	db        *sql.DB
}

func NewCrawler(queue *queue.MessageQueue, db *sql.DB, limit *colly.LimitRule, opts ...colly.CollectorOption) (*Crawler, error) {
	c := colly.NewCollector(
		opts...,
	)

	err := c.Limit(limit)
	if err != nil {
		return nil, fmt.Errorf("couldn't configure limit: %w", err)
	}

	return &Crawler{
		collector: c,
		queue:     queue,
		db:        db,
	}, nil
}

func (c *Crawler) ConsumeTasks(messageQueue *queue.MessageQueue) {
	msgs, err := messageQueue.GetConsumer()
	if err != nil {
		log.Fatalf("couldn't consume from queue: %s", err)
	}

	go func() {
		var respCode int
		log.Printf("Processing...")

		for d := range msgs {
			if d.ContentType != "application/json" {
				continue
			}

			var articles []*crawlerDb.Article
			var authorHyperlinks []string

			task := &Task{}
			err := json.Unmarshal(d.Body, task)
			if err != nil {
				log.Fatalf("Couldn't unmarshal task: %s", err)
			}

			u, err := url.Parse(task.Url)
			if err != nil {
				log.Fatalf("Couldn't parse url: %s", err)
			}

			switch v := parsePath(u.Path); v {
			case IndexPage:
				log.Println("Processing page index task")
				authorHyperlinks, respCode, err = c.CrawlIndexPageOfAuthors(task.Url)
				for i, s := range authorHyperlinks {
					if i > limitHyperlinkPerPage {
						break
					}

					publishErr := c.PublishCrawlTak(messageQueue, &Task{
						Url: s,
					})
					if publishErr != nil {
						log.Fatalf("Failed to publish author page crawl task: %v", err)
					}
				}

			case AuthorPage:
				log.Println("Processing extracting articles task")
				articles, respCode, err = c.CrawlPageOfAuthor(task.Url)
				log.Printf("Scrapped articles: %v\n", len(articles))

			default:
				log.Printf("Unexpected path given :%v", v)
			}

			if err != nil {
				log.Printf("Couldn't crawl url: %s, %v", task.Url, err)

				if respCode == http.StatusGatewayTimeout {
					log.Println("CRAWLER GOT TIMED OUT")

					err := d.Reject(true)
					if err != nil {
						log.Printf("Failed to NOT acknowledge message: %s", err)
					}

					log.Printf("Sleeping %v seconds before starting crawling tasks again", sleepTimeoutDuration)
					time.Sleep(sleepTimeoutDuration)
					respCode = 0
				}

				continue
			}

			if err := d.Ack(false); err != nil {
				log.Printf("Failed to acknowledge message: %s", err)
				continue
			}
			log.Printf("[x] Ack task for %s", task.Url)

			if len(articles) > 0 {
				err = crawlerDb.BulkInsertArticles(c.db, articles)
				if err != nil {
					log.Printf("Failed to bulk insert articles: %s", err)
				}
			}

			time.Sleep(poolingTimeoutDuration)
		}
	}()
}

func parsePath(path string) UrlPath {
	split := strings.Split(path, "/")
	if len(split) == 1 {
		return UnkownPage
	}
	split = split[1:]
	fmt.Println(split, len(split))

	if len(split) == 1 {
		return IndexPage
	}

	return AuthorPage
}

func (c *Crawler) CrawlPageOfAuthor(url string) ([]*crawlerDb.Article, int, error) {
	var articles []*crawlerDb.Article
	var gotErr error
	var errorStatusCode int

	c.collector.OnHTML("cite.data.tts-content[itemprop=\"headline\"]", func(e *colly.HTMLElement) {
		articleTitle := e.ChildText("span.title")
		authors := e.ChildTexts("span[itemprop=\"author\"]")

		articles = append(articles, &crawlerDb.Article{
			Title:   strings.Join(strings.Fields(strings.TrimSpace(articleTitle)), " "),
			Authors: authors,
		})

	})

	c.collector.OnError(func(resp *colly.Response, err error) {
		gotErr = err
		errorStatusCode = resp.StatusCode
	})

	gotErr = c.collector.Visit(url)
	if gotErr != nil {
		return nil, errorStatusCode, gotErr
	}
	log.Println("Initial visit", url)

	return articles, errorStatusCode, gotErr
}

func (c *Crawler) CrawlIndexPageOfAuthors(url string) ([]string, int, error) {
	var authorHyperLinks []string
	var gotErr error
	var errorStatusCode int

	c.collector.OnHTML("div.columns.hide-body", func(e *colly.HTMLElement) {
		authorHyperLinks = e.ChildAttrs("a", "href")
	})

	c.collector.OnError(func(resp *colly.Response, err error) {
		gotErr = err
		errorStatusCode = resp.StatusCode
	})

	gotErr = c.collector.Visit(url)
	if gotErr != nil {
		return nil, errorStatusCode, gotErr
	}
	log.Println("Initial visit", url)

	return authorHyperLinks, errorStatusCode, gotErr
}

func (c *Crawler) PublishCrawlTak(messageQueue *queue.MessageQueue, task *Task) error {
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
			return nil
		}
		log.Printf("Failed to publish message: %v", err)
	}

	return nil
}
