Article Crawler with Docker

This project is a lightweight web crawler designed to search for and collect articles from specified sources. Built using Go, RabbitMQ, and PostgreSQL, it provides an easy-to-use API endpoint to trigger the crawling process. The application is fully containerized using Docker, making deployment and setup simple and seamless.
Features

    Article Crawling: Automatically crawls a source website to search for articles.
    Asynchronous Processing: Utilizes RabbitMQ to handle crawling tasks efficiently.
    Database Integration: Stores collected articles in a PostgreSQL database for easy access and management.
    REST API: Provides an endpoint (/startCrawl) to initiate the crawling process.

Prerequisites

Before running this project, ensure you have the following installed on your system:

    Docker: Used to containerize and run the application and its dependencies.
    Docker Compose: Simplifies running the multi-container application with a single command.

Getting Started
1. Clone the Repository

git clone https://github.com/yourusername/article-crawler.git
cd article-crawler

2. Start the Application

Run the following command to build and start the services:

docker-compose up

This command will:

    Start the RabbitMQ service for managing crawling tasks.
    Start the PostgreSQL database for storing crawled articles.
    Start the Go-based API for controlling the crawler.

3. Access the API

Once the services are up, open your browser or use a tool like curl to visit the following URL:

http://localhost:8081/startCrawl

This endpoint will initiate the crawling process. The crawler will start searching for articles and store the results in the PostgreSQL database.
Project Structure

    Go Application: Contains the main crawler logic and the REST API.
    RabbitMQ: Handles task distribution for asynchronous crawling.
    PostgreSQL: Stores the articles, including metadata like titles and authors.

Example Workflow

    Start the Docker containers using docker-compose up.
    Trigger the crawler via the API (/startCrawl).
    Articles are fetched from the source and stored in the database.
    Use additional APIs (if implemented) to retrieve or analyze the stored articles.

Environment Variables

The application uses environment variables for configuration. You can set these in the .env file or pass them directly when running docker-compose. Default values are provided for development convenience.

Key variables include:

    DB_HOST, DB_PORT, DB_USER, DB_PASS, DB_NAME for database connection.
    RABBITMQ_HOST, RABBITMQ_PORT, AMPQ_USER, AMPQ_PASS for RabbitMQ.
