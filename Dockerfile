FROM golang:1.25 as builder
WORKDIR /build

COPY go.mod .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bot cmd/bot/main.go

FROM scratch
WORKDIR /app

COPY --from=builder /build/configs/values_ci.yaml .
COPY --from=builder /build/bot .

ENV CONFIG_FILE=values_ci.yaml

EXPOSE 8080

ENTRYPOINT [ "./bot" ]
