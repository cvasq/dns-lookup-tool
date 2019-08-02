APP_NAME=dns-lookup-tool
FRONTEND_DIR=ui

all: build-vue-app go docker 
build: frontend go-build
frontend: build-vue-app
go-build: go
container: docker

export VUE_APP_BASE_PATH=/
export VUE_APP_WS_URL=ws://localhost:8080/dns-check

.PHONY: build-vue-app
build-vue-app: 
	npm --prefix $(FRONTEND_DIR) install
	npm run --prefix $(FRONTEND_DIR) build

.PHONY: go
go:
	go generate .
	go build .

.PHONY: docker
docker:
	docker build -t cvasquez/$(APP_NAME) .
