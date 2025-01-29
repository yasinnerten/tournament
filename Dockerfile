FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy
RUN go mod vendor

COPY . .

# Copy the .env file for debugging purposes
COPY .env .env
# List the contents of the /app directory to verify the presence of the .env file
RUN ls -la /app/scripts

# Run the build script
RUN /app/scripts/build.sh

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/bin/app .

# Copy the .env file for debugging purposes
COPY .env .env

# Copy the generated Swagger documentation
COPY docs /app/docs

EXPOSE 8080

CMD ["./app"]