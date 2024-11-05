# This Docker boilerplate does not provide shell access, primarily because this repository does not depend on the operating system.
# Instead, it relies solely on the server's resources, such as CPU, memory, network, and disk (for serving static files for the front-end site).
#
# Use the official Golang image to create a build artifact.
# This is a multi-stage build. This stage is named 'builder'.
#
# Note: This is a boilerplate for Docker regarding this repository.
# Also, note that the "# Copy the source code." or other "# Copy ..." comments need to be written specifically for your use case,
# for example, copying a directory of the source code for building a container.
FROM golang:1.23.2 AS builder

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /app

# Copy the go.mod and go.sum files to download all dependencies.
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download

# Copy the backend source code.


# Additionally, copy the frontend source code (e.g., root tailwind.config.js, /frontend/public/magic_embedded.go, htmx templ).


# Build the application.
#
# Note: This design might require experimental C + Go, so "-installsuffix cgo" won't work anyway due to CGO_ENABLED=0. If CGO_ENABLED=1, it would work.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /restapis ./backend/cmd/server/run.go

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/
#
# Note: Consider using a minimalist image that does not rely on the operating system. The only dependency needed for this repository is the ca-certificates package,
# in case you need to make calls to HTTPS endpoints.
FROM alpine:latest

# Install the ca-certificates alpine package in case you need to make calls to HTTPS endpoints.
# The scratch image is the most minimal image in Docker. This image is only 5MB and has no shell.
# If you need to debug within the container, you might want to use a different base image.
#
# Also, note that while this original Dockerfile boilerplate is free from maintaining the image because it's zero-vulnerable (for example, no CVEs when scanning),
# the only thing to focus on is writing the code for this repository.
RUN apk --no-cache add ca-certificates

# This is a safe way
WORKDIR /root/

# Copy the pre-built binary file from the previous stage.
COPY --from=builder /restapis .

# Note: This is where the magic embedding (see /frontend/public/magic_embedded.go) takes place after the builder stage source code is done.
#
# Copy the docs directory (e.g., contains the openapi.json file) to the image.

# Note: This is where the magic embedding (see /frontend/public/magic_embedded.go) takes place after the builder stage source code is done.
#
# Copy the frontend assets (e.g., js, css, other static files such as png, jpeg) to the image.


# Expose port 8080 to the outside world.
# This can be modified.
EXPOSE 8080

# Command to run the executable.
CMD ["./restapis"]
