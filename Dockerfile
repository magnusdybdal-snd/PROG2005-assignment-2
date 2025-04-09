FROM golang:1.24.1 AS builder

LABEL maintainer="magndy@stud.ntnu.no"
LABEL stage=builder

# Set up execution environment in container's GOPATH
WORKDIR /go/src/app/

# Copy relevant folders into container
COPY ./main.go /go/src/app/main.go
COPY ./handlers /go/src/app/handlers
COPY ./utils /go/src/app/utils
COPY ./go.mod /go/src/app/go.mod
COPY ./go.sum /go/src/app/go.sum

# Compile binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o server

# Indicate port on which server listens
EXPOSE 8080

# Instantiate binary
CMD ["./server"]
