.PHONY: wire swag

GOCMD=go

# format to write make file is 
# target: prerequesttobeececuted 1 if any, prerequesttobeececuted 2 if any,..
# 	actions...

wire:
	cd pkg/di && wire

swag:
	swag init -g pkg/api/handler/user.go -o ./cmd/api/docs

run:
	@echo "Smart_Store Server running...."  
#	go run cmd/api/main.go
#	go run ./cmd/api
	$(GOCMD) run ./cmd/api
	





