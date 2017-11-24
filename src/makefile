folder ?= .
bin ?= 

help:
	@echo "Go distributed system"
	@echo ""
	@echo " dependencies   gets all the required dependencies"
	@echo " build        compile and installs the binary"
	@echo " install      get all the dependencies and install dependencies"

dependencies:

build: 
	@go build -a -o $(bin)
	@cp $(bin) $(folder)/$(bin)

install: dependencies build

.PHONY: help dependencies build install
