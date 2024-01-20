# Use an official Go runtime as a parent image
FROM golang:1.18

# Set the working directory inside the container
WORKDIR /go/src/erp

# Install air for live code reloading
# RUN go install github.com/cosmtrek/air@latest
RUN apt-get -y update; apt-get -y install curl
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Copy the Go module files and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the current directory contents into the container at /go/src/app
COPY . .

# Build the Go application
RUN go build -o erp

# Expose port 8080 for the application
EXPOSE 8080

# Run the application with air for live reloading
CMD air
