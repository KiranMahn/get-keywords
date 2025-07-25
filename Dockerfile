# Use a multi-stage build to reduce the final image size

# To run this Dockerfile, ensure you have Docker installed and run:
# docker build -t get-keywords . && docker run -p 8081:8081 get-keywords

# Stage 1: Build the Go application
FROM golang:1.24-bullseye AS builder

# Set the working directory inside the container
WORKDIR /go/src

# Copy the Go modules files and download dependencies
COPY ./ ./

# Build the Go application
RUN set -x && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o serve *.go

# Stage 2: Create the final image
FROM golang:1.24-bullseye
WORKDIR /app

# Copy the built application from the builder stage
COPY --from=builder /go/src/ ./
COPY --from=builder /go/src/serve ./

EXPOSE 8081
# Command to run the executable
CMD [ "/app/serve" ]