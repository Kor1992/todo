FROM golang:1.26.1-bookworm AS builder
ENV GOPROXY=https://goproxy.io,direct
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/cmd/todoapp/exe /app/cmd/todoapp

FROM alpine:3.23
WORKDIR /app
COPY --from=builder /app/cmd/todoapp/exe /app
COPY --from=builder /app/docs ./docs   
CMD ["/app/exe"]
COPY --from=builder /app/docs ./docs