# This Docker boilerplate does not provide shell access, primarily because this repository does not depend on the operating system.
# Instead, it relies solely on the server's resources, such as CPU, memory, network, and disk (for serving static files for the front-end site).
#
# Use the official Golang image to create a build artifact.
# This is a multi-stage build. This stage is named 'builder'.
#
# Note: This is a boilerplate for Docker regarding this repository.
# Also, note that the "# Copy the source code." or other "# Copy ..." comments need to be written specifically for your use case,
# for example, copying a directory of the source code for building a container.
#
# This uses a custom base image on Alpine:latest because the official Docker Golang images can be slow the maintainer to update with new versions.
#
# Repo: https://git.b0zal.io/H0llyW00dzZ/golang.git
#
# TODO: Automate updates when new versions are available, instead of waiting for the slow maintainer.
#
# Known Bug: Gitea manifests don't work with versions v1.23.3 and 1.23.3.
# failed commit on ref "index-sha256:524acd083758062a04371a283008aa7d6a82d678fbf479c29a45f5ba86a04c57": unexpected status from PUT request
#
# Additionally, this has been modified (to previously) to support multiple build architectures.
FROM golang:1.23.4 AS builder

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
# Additionally, if the TAG variable is empty from Variables Management, the build will proceed without any tags (default Go build).
# Ensure tags are correct, as it is not possible to handle errors in this Dockerfile.
ARG TAG
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags=${TAG} -o /restapis ./backend/cmd/server/run.go

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


# Label for improved versioning
#
# TODO: Add a description based on the latest commits (e.g., 100+ commits), if possible, without using "run" commands in CI/CD.
# Using "run" commands for this can indicate poor GitOps, DevOps, and DevSecOps practices.
ARG VENDOR
ARG REPO
ARG VERSION
ARG TARGETPLATFORM
ARG SHA

LABEL vendor="${VENDOR}"
LABEL repo="${REPO}"
LABEL version="${VERSION}"
LABEL platform="${TARGETPLATFORM}"
LABEL sha="${SHA}"

# Expose port 8080 to the outside world.
# This can be modified.
EXPOSE 8080

# Command to run the executable.
CMD ["./restapis"]
