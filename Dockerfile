# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

RUN git config --global url."git@github.com:".insteadOf "https://github.com"
RUN git config --global url."https://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:x-oauth-basic@github.com/".insteadOf "https://github.com/"


# Copy go mod and sum files
COPY go.mod go.sum ./

# Copy the source from the current directory to the Working Directory inside the container
ADD . /app

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -o main .

# Command to run the executable
CMD ["./main"]