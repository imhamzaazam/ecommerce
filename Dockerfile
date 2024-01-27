FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./codegen ./codegen

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./transactions/server


CMD ["/app/server"]

