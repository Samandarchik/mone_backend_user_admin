# Dockerfile
FROM golang:1.24-alpine AS builder

# Kerakli paketlarni o'rnatish
RUN apk add --no-cache git gcc musl-dev

# Ish katalogini yaratish
WORKDIR /app

# Go modullarini nusxalash
COPY go.mod go.sum ./

# Go modullarni yuklab olish
RUN go mod download

# Barcha kodlarni nusxalash
COPY . .

# Binaryni qurish
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage - minimal image
FROM alpine:latest

# SSL sertifikatlarini qo'shish (HTTPS so'rovlar uchun)
RUN apk --no-cache add ca-certificates wget

# Ish katalogi
WORKDIR /root/

# Binary va kerakli kataloglarni nusxalash
COPY --from=builder /app/main .

# Data katalogini nusxalash (agar kerak bo'lsa)
RUN mkdir -p data uploads

# Ma'lumotlar uchun volume
VOLUME ["/root/data", "/root/uploads"]

# Portni ochish
EXPOSE 1010

# Serverni ishga tushirish
CMD ["./main"]