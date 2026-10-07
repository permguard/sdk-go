# Copyright (c) 2022 Nitro Agility S.r.l.
# SPDX-License-Identifier: Apache-2.0

.DEFAULT_GOAL := build

brew:
	brew install golangci-lint
	brew install staticcheck
	brew install gofumpt
	brew install protobuf

install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.5
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	go install github.com/google/addlicense@latest

clean:
	rm -rf dist/
	rm -rf tmp/
	rm -f coverage.out
	rm -f result.json

init-dependency:
	go get -u google.golang.org/grpc v1.69.4
	go get -u google.golang.org/protobuf v1.36.3

mod:
	go mod download
	go mod tidy

protoc:
	protoc proto/v1/*.proto \
		--go_out=internal/grpc/v1 --go_opt=paths=source_relative \
		--go-grpc_out=internal/grpc/v1 --go-grpc_opt=require_unimplemented_servers=false,paths=source_relative \
		--proto_path=proto/v1
	bash scripts/header-generated.sh internal/grpc/v1/pdp.pb.go internal/grpc/v1/pdp_grpc.pb.go

check:
	staticcheck  ./...

lint:
	go vet ./...
	gofmt -s -w **/**.go
	gofumpt -l -w .
	golangci-lint run --disable-all --enable staticcheck

lint-fix:
	gofmt -s -w **/**.go
	go vet ./...
	gofumpt -l -w .
	golangci-lint run ./... --fix

test:
	go test ./...

teste2e:
	export E2E="TRUE" && GOFLAGS="-count=1" go test ./e2e/...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

coverage-plugin:
	go test -coverprofile=coverage.out ./plugin/...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

converage-%:
	go test -coverprofile=coverage.out ./...

converage-json:
	go test -json -coverprofile=coverage.out ./... > result.json

# disallow any parallelism (-j) for Make. This is necessary since some
# commands during the build process create temporary files that collide
# under parallel conditions.
.NOTPARALLEL:

.PHONY: clean mod lint lint-fix
