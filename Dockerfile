FROM golang:1.15-alpine AS build

RUN apk add --update git
ENV GO111MODULE=on CGO_ENABLED=0
WORKDIR /code/

# Precache gomod dependencies
COPY go.mod go.sum /code/
RUN go mod download

COPY . /code/
RUN go build -o /bin/cloudforecast-agent

FROM alpine:latest
COPY --from=build /bin/cloudforecast-agent /bin/cloudforecast-agent
ENTRYPOINT ["/bin/cloudforecast-agent"]
