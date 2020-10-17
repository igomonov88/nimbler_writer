SHELL := /bin/bash

all:

proto:
	protoc --go_out=plugins=grpc:. nimbler_reader/proto/contract.proto
	protoc --go_out=plugins=grpc:. nimbler_writer/proto/contract.proto


migrate:
	go run ./cmd/admin/main.go migrate

seed: migrate
	go run ./cmd/admin/main.go seed

writer-api:
	docker build \
		-f users/dockerfile.nimbler_writer \
		-t igorgomonov/nimbler_writer-amd64:1.0 \
		--build-arg PACKAGE_NAME=nimbler_writer \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +”%Y-%m-%dT%H:%M:%SZ”` \
		.

metrics:
	docker build \
		-f sugar/dockerfile.metrics \
		-t igorgomonov/nimbler_metrics-amd64:1.0 \
		--build-arg PACKAGE_NAME=metrics \
		--build-arg PACKAGE_PREFIX=sidecar/ \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +”%Y-%m-%dT%H:%M:%SZ”` \
		.

up:
	docker-compose up

down:
	docker-compose down

test:
	go test -mod=vendor ./... -count=1

clean:
	docker system prune -f

stop-all:
	docker stop $(docker ps -aq)

remove-all:
	docker rm $(docker ps -aq)

tidy:
	go mod tidy
	go mod vendor

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -t -d -v ./...

deps-cleancache:
	go clean -modcac