 Bag of Crawler Tasks

---

## Features

- **Article Crawling**: Automatically crawls a source website to search for articles.
- **Asynchronous Processing**: Utilizes RabbitMQ to handle crawling tasks efficiently.
- **Database Integration**: Stores collected articles in a PostgreSQL database for easy access and management.
- **REST API**: Provides an endpoint (`/startCrawl`) to initiate the crawling process.

---

## Prerequisites

Before running this project, ensure you have the following installed on your system:

- **Docker**: Used to containerize and run the application and its dependencies.
- **Docker Compose**: Simplifies running the multi-container application with a single command.

---

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/article-crawler.git
cd article-crawler

docker-compose up -d
```


### 2. Acess the API

http://localhost:8081/startCrawl

### 3. Look into docker container logs to see the ongoing activity
