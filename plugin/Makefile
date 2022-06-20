.DEFAULT_GOAL := build

PKG_NAME = vault-plugin-auth-google-acc
BUILD_FLAGS =

build:
	go build $(BUILD_FLAGS) -o bin/$(PKG_NAME) cmd/$(PKG_NAME)/main.go

release: BUILD_FLAGS += -ldflags "-w -s"
release: build

format:
	gofmt -s -w .
