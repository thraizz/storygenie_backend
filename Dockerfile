# Start from the official golang image
FROM golang:latest

# Set the working directory
WORKDIR /app

# Copy the source code to the container
COPY . .

# Build the Go application
RUN go build -o storygenie-backend

# Expose port 8080 to the outside world
EXPOSE 8080

# Start the Go application
CMD ["./storygenie-backend"]
