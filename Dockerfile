FROM golang:1.26.0-alpine AS builder

WORKDIR /loyalty-card-service

RUN apk --no-cache add git bash make gcc gettext musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

ENV CONFIG_PATH=config/config.yaml
ENV CGO_ENABLED=0

RUN go build --ldflags="-w -s" -o ./build/loyalty-card-service ./cmd/loyalty-card-service

FROM alpine AS runner

RUN apk add --no-cache ca-certificates

WORKDIR /loyalty-card-service

COPY --from=builder /loyalty-card-service/config /loyalty-card-service/config
COPY --from=builder /loyalty-card-service/build/loyalty-card-service /loyalty-card-service/build/loyalty-card-service

EXPOSE 8080

CMD ["./build/loyalty-card-service"]