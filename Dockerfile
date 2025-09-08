# Dockerfile
FROM golang:1.23-alpine AS builder

# Kerakli paketlarni o'rnatish
RUN apk add --no-cache git

# Ish katalogini yaratish
WORKDIR /app

# Go modullarini nusxalash va yuklab olish
COPY go.mod go.sum ./
RUN go mod download

# Barcha kodlarni nusxalash
COPY . .

# Binaryni qurish
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main main.go models.go data_service.go middleware.go printer_service.go routes.go

# Final stage - minimal image
FROM alpine:latest

# SSL sertifikatlarini qo'shish (HTTPS so'rovlar uchun)
RUN apk --no-cache add ca-certificates

# Ish katalogi
WORKDIR /root/

# Binary va data katalogini nusxalash
COPY --from=builder /app/main .
COPY --from=builder /app/data ./data

# Ma'lumotlar uchun volume
VOLUME ["/root/data"]

# Portni ochish
EXPOSE 1010

# Serverni ishga tushirish
CMD ["./main"]