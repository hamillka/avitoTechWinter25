APP=avito-shop-service
.PHONY: linter

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
	go test ./e2e_tests -count=1 -coverpkg=./...
	make stop

all_tests:
	make units
	make e2e-w-coverage

swag-gen:
	swag init -g ../../cmd/avito-shop-service/main.go -o ./api -d ./internal/handlers
