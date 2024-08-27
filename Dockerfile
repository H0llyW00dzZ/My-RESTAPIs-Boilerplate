# This Docker boilerplate has no shell access, mostly because this repository does not rely on the operating system,
# only on the resources of the machine server such as CPU, memory, network and disk.
#
# Use the official Golang image to create a build artifact.
# This is a multi-stage build. This stage is named 'builder'.
#
# Note: This is a boilerplate for Docker regarding this repository.
# Also, note that the "# Copy the source code." or other "# Copy ..." comments need to be written specifically for your use case,
# for example, copying a directory of the source code for building a container.
FROM golang:1.23.0 AS builder

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /app

# Copy the go.mod and go.sum files to download all dependencies.
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download

# Copy the source code.


# Additionally, copy the frontend assets


# Build the application.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /restapis ./backend/cmd/server/run.go

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/
FROM alpine:latest

# Install ca-certificates in case you need to make calls to HTTPS endpoints.
# The scratch image is the most minimal image in Docker. This image is only 5MB and has no shell.
# If you need to debug within the container, you might want to use a different base image.
RUN apk --no-cache add ca-certificates

# Which is safe-way
WORKDIR /root/

# Copy the pre-built binary file from the previous stage.
COPY --from=builder /restapis .

# Copy the docs directory (which contains the openapi.json file) to the image.


# Copy the frontend assets to the image.


# Expose port 8080 to the outside world.
# This can be modified
EXPOSE 8080

# Command to run the executable.
CMD ["./restapis"]
