# Use an official Go runtime as a parent image
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /go/src/erp

# Install air for live code reloading
# RUN go install github.com/cosmtrek/air@latest
RUN apt-get -y update; apt-get -y install curl
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

RUN apt-get update && apt-get install -y git

# Copy the Go module files and download dependencies
COPY go.mod .
# COPY go.sum .
RUN go mod download

# Exclude .git and .gitignore
COPY .dockerignore .dockerignore

# Copy the current directory contents into the container at /go/src/app
COPY . .

# Build the Go application
# RUN go build -o erp -ldflags="-w -s" -tags="netgo" -installsuffix=netgo -buildmode=pie -buildvcs=false

RUN if [ "$VCS" = "off" ]; then \
    go build -o erp -ldflags="-w -s" -tags="netgo" -installsuffix=netgo -buildmode=pie; \
  else \
    go build -o erp; \
  fi

# Expose port 8080 for the application
EXPOSE 8080

# Run the application with air for live reloading
CMD air