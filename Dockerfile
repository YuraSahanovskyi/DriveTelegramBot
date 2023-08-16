FROM golang:1.20 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/bot .
COPY ./pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux go build -o /bot-bin .

FROM build-stage AS test-stage
RUN go test -v ./...

FROM alpine:3.18 AS build-release-stage

WORKDIR /app

COPY --from=build-stage /bot-bin ./bot-bin
COPY ./config ./config


EXPOSE 8080

CMD [ "./bot-bin" ]