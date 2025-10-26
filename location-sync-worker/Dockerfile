FROM golang:alpine AS builder

WORKDIR /app
# RUN apk add --no-cache --upgrade git openssh-client make

COPY . .

RUN go env -w GO111MODULE=on
RUN go mod download
RUN go generate ./...

RUN GOOS=linux go build -o main ./

## Distribution Development
FROM alpine AS final

WORKDIR /app

COPY --from=builder /app/main ./
COPY --from=builder /app/.env.example ./

CMD [ "/app/main" ]