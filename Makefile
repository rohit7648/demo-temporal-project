GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

ifeq ($(GOHOSTOS), windows)
	# The `find.exe` is different from `find` in bash/shell.
	# Refer to https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	# Use git-bash.exe to run find cli or other cli-friendly tools, since every developer has Git.
	# Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find api-protos/feature-store -name *.proto)
endif

.PHONY:	init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest

.PHONY: api
# generate api proto
api:
	protoc --proto_path=api-protos/demo_temporal_project/v1 \
			--proto_path=api-protos/third_party \
			--go_out=paths=source_relative:api-protos/demo_temporal_project/v1 \
			--go-http_out=paths=source_relative:api-protos/demo_temporal_project/v1 \
			--go-grpc_out=paths=source_relative:api-protos/demo_temporal_project/v1 \
			--go_opt=Mtransaction.proto=demo_temporal_project/api-protos/demo_temporal_project/v1 \
			 --go-http_opt=Mtransaction.proto=demo_temporal_project/api-protos/demo_temporal_project/v1 \
			 --go-grpc_opt=Mtransaction.proto=demo_temporal_project/api-protos/demo_temporal_project/v1 \
			transaction.proto


.PHONY:	config
# generate internal proto
config:
	protoc --proto_path=./internal \
	       --proto_path=./third_party \
	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY:	build
# build
build:	
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY:	generate
# generate
generate:
	go mod tidy
	go get github.com/google/wire/cmd/wire@latest
	go generate ./...

.PHONY:	all
# generate all
all:
	make generate;

run:
	make build;
	./bin/application;

run-worker:
	make build;
	./bin/worker;

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

migrate-create:  ### create new migration
	migrate create -ext sql -dir db/migrations '$(title)'
.PHONY: migrate-create

migrate-up: ### migration up
	migrate -path db/migrations -database '$(PG_URL)?sslmode=disable' up
.PHONY: migrate-up

migrate-down: ### migration down
	migrate -path db/migrations -database '$(PG_URL)?sslmode=disable' down
.PHONY: migrate-down

