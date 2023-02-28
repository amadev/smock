FROM golang:alpine

WORKDIR /app

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /app/build && chmod +x /app/build

RUN ls -l /app

EXPOSE 8080

CMD [ "/app/build" ]
