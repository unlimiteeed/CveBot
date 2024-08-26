# Use the official Golang image with the latest stable version
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o cve-notifier cmd/cveBot/cveBot.go

# Run the Go application
CMD ["./cve-notifier"]
