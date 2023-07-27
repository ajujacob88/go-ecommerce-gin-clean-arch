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


docker-build: ## To build new docker image #this is only required to build a docker image, remember while at the time of buiding docker image(running this command), dont forget to change the .env file's DB-HOST to db instead of local host...DB_HOST=db instead of DB_HOST=localhost.. at the time of running this command, otherwise the docker image wont work
	docker build -t ajujacob/smarstore-ecommerce-api:0.0.1.RELEASE .    
# . means the docker file is in the current directory	

docker-up: ## To up the docker compose file
	docker-compose up 

docker-down: ## To down the docker compose file
	docker-compose down

