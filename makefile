.PHONY: wire swag

# format to write make file is 
# target: prerequesttobeececuted 1 if any, prerequesttobeececuted 2 if any,..
# 	actions...

wire:
	cd pkg/di && wire

swag:
	swag init -g pkg/api/handler/user.go -o ./cmd/api/docs

run:
	go run cmd/api/main.go 





