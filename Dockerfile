# Dockerfile

FROM golang:1.14

# Set the Work Directory
WORKDIR /usr/src/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build
EXPOSE 8080

RUN chmod u+x server.sh
CMD ["./server.sh"]