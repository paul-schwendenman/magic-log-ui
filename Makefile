# Makefile

NAME := magic-log
VERSION ?= dev
BUILD_DIR := cmd
WEB_DIR := web
STATIC_DIR := $(BUILD_DIR)/static

.PHONY: all clean build frontend backend run test coverage

all: build

clean:
	rm -rf $(STATIC_DIR) $(NAME)

frontend:
	cd $(WEB_DIR) && pnpm install && pnpm build
	mkdir -p $(STATIC_DIR)
	cp -r $(WEB_DIR)/build/* $(STATIC_DIR)/

backend:
	go build -ldflags "-X main.version=$(VERSION)" -o $(NAME) ./$(BUILD_DIR)

build: frontend backend

run: build
	./$(NAME)

test:
	go test ./... -v

coverage:
	go test ./internal/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report written to coverage.html"

version:
	@echo $(VERSION)
