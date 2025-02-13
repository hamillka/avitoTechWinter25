APP=avito-shop-service

build:
	docker-compose build $(APP)

run:
	docker-compose up -d $(APP)

stop:
	docker-compose down
