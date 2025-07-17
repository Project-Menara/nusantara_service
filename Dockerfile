# FROM golang:1.23.1-alpine

# WORKDIR /app

# # COPY file mod ke /app
# COPY go.mod go.sum ./

# # Install dependency
# RUN go mod tidy

# # Copy semua isi project ke container
# COPY . .

# # Build binary
# RUN go build -o main .

# # Expose port (jika kamu pakai port 8080 misalnya)
# EXPOSE 8080

# # Jalankan aplikasi
# CMD ["./main"]

FROM golang:1.23.1-alpine

WORKDIR /app

ARG CACHEBUST=1

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
