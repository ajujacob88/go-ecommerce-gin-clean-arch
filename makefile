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


test: ## Run tests
# go test ./... -v -cover
	$(GOCMD) test ./... -v -cover


mockgen: ## Generate mock repository and usecase functions	
	mockgen -source=pkg/repository/interface/admin.go -destination=pkg/mock/repositoryMock/adminMock.go -package=repositoryMock
	mockgen -source=pkg/usecase/interface/admin.go -destination=pkg/mock/usecaseMock/adminMock.go -package=usecaseMock



