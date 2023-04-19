# Use the official Go image as the base image
FROM golang:1.19

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY go.mod ./
COPY go.sum ./

# Download the Go dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go application
RUN go build -o websocket-server

# Expose the port the server will listen on
EXPOSE 8080

# Start the server
CMD ["/app/websocket-server"]
