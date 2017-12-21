.PHONY: help dependencies tests build install start

bin = twitter-votes
folder = ./bin
image = jaumecapdevila/twitter-votes

help:
	@echo "Twitter votes"
	@echo ""
	@echo " dependencies gets all the required dependencies"
	@echo " build compile the go binary"
	@echo " install runs both previous steps"

dependencies:
	@echo "Installing package dependencies"
	@go get -u github.com/kardianos/govendor
	@govendor sync

tests:
	@echo "Executing unit tests..."
	@govendor test +local

build: dependencies
	@echo "Building the $(bin) binary"
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(folder)/$(bin)
	@echo "Building docker image"
	docker build -t $(image) .

setup-env:
	@echo "Exporting all the needed env variables..."
	@chmod +x setup.sh
	@./setup.sh

start: dependencies build
	@echo "Starting twitter votes service"
	cd ../
	docker-compose up -d app
