FROM golang:1.23.1-alpine

WORKDIR /app

# Tambahkan arg agar cache invalidated
ARG CACHE_BUST=1

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

# Ini build ulang binary Go
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
