APP_NAME=dns-lookup-tool
FRONTEND_DIR=frontend

all: build-vue-app build-go docker
frontend: build-vue-app
go-app: build-go

.PHONY: build-vue-app
build-vue-app: 
	npm run --prefix $(FRONTEND_DIR) build

.PHONY: build-go
build-go:
	go build .

.PHONY: docker
docker:
	docker build -t cvasquez/$(APP_NAME) .
