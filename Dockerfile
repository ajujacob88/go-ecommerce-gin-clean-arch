#reference - https://docs.docker.com/language/golang/build-images/
# Base image : specify the enviornment
FROM golang:1.20-alpine AS build-stage

# Maintainer info
LABEL maintainer="Aju Jacob <ajujacob88@gmail.com>"

# working directory
WORKDIR /home/app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -v -o /home/build/api ./cmd/api


FROM gcr.io/distroless/static-debian11 AS build-release-stage

# Maintainer info
LABEL maintainer="Aju Jacob <ajujacob88@gmail.com>"

COPY --from=build-stage /home/build/api /api
COPY --from=build-stage /home/app/views /views
COPY .env /

EXPOSE 3000

CMD ["/api"]
