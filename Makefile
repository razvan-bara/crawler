restart:
	docker-compose down && docker-compose up -d --build

start-test-env:
	docker-compose down && IS_TST=1 docker-compose up -d --build