APP=avito-shop-service

build:
	docker-compose build $(APP)

run:
	docker-compose up -d $(APP)

stop:
	docker-compose down

linter:
	golangci-lint run ./... --config=./.golangci.yaml

coverage:
	go test ./... -v -coverpkg=./...
