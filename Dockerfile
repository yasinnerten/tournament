# Use the official Golang image as a build stage
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

# Build the Go application
RUN go build -o main .

# Use the official Golang image for the final executable
FROM golang:1.23

WORKDIR /app

# Copy the source code and built executable from the builder stage
COPY --from=builder /app .

# Ensure the executable has the correct permissions
RUN chmod +x ./main

# Command to run the application
CMD ["go", "run", "main.go"]