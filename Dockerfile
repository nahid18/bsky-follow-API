# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory in the container
WORKDIR /app

# Copy the Go REST API server code and the CLI binary into the container
COPY main.go .
COPY bsky .
COPY go.mod .
COPY go.sum .

# Build the Go REST API server
RUN go build -o api main.go

# Expose the port on which the API server will listen
EXPOSE 8080

# Define the command to run when the container starts
CMD ["./api"]
