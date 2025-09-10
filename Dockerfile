FROM golang:1.22.3-alpine AS builder

LABEL authors="jexroid"

WORKDIR /build
COPY .env .
COPY . .

RUN export GO111MODULE=on

RUN go mod download

RUN go test -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o ./gopi

# Runnng Go Application in NonRoot Debian User
FROM debian:latest

# Update the package list, install sudo, create a non-root user, and grant password-less sudo permissions
# its recomended to use arguments for UserID or GroupID but for ore simplified purses it's static
RUN addgroup --gid 51 nonroot && \
    adduser --uid 14 --gid 51 --disabled-password --gecos "" nonroot && \
    echo 'nonroot ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers

USER nonroot

WORKDIR /home/nonroot/app

RUN chmod -R 755 /home/nonroot/app

COPY --from=builder /build/gopi /home/nonroot/app/gopi
COPY .env /home/nonroot/app

EXPOSE 8000

CMD ["/home/nonroot/app/gopi"]