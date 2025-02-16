APP=avito-shop-service
.PHONY: linter build run stop linter units units-coverage e2e-w-coverage all_tests swag-gen load_testing

build:
	docker-compose build $(APP)

run:
	docker-compose up -d $(APP)

stop:
	docker-compose down

linter:
	golangci-lint run ./... --config=./.golangci.yaml

units:
	go test ./internal/services/unit_tests/
	go test ./internal/repositories/unit_tests/

units-coverage:
	go test ./... -coverpkg=./...

e2e-w-coverage:
	make run
	sleep 5
	go test ./e2e_tests -coverpkg=./...
	make stop

all_tests:
	make units
	make e2e-w-coverage

swag-gen:
	swag init -g ../../cmd/avito-shop-service/main.go -o ./api -d ./internal/handlers

load_testing:
	make run
	go run ./load_testing/db_100k_users.go
	k6 run ./load_testing/load.js
	make stop
