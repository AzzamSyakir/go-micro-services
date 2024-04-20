# Start from golang base image.
FROM golang:1.22-alpine

# Install dependencies.
RUN apk update
RUN apk add --no-cache git build-base

# Setup workdir.
ENV WORKDIR /app
WORKDIR $WORKDIR

# Copy sources from the host current directory to the workdir inside the container.
COPY . .


# Install the modules.
RUN go install -v ./...

# Clean all the downloaded modules.
RUN go clean -modcache

# Build the Go app.
RUN go build -C $WORKDIR/cmd/web -o $WORKDIR/build

# Copy the build to the final image.
FROM alpine:latest

# Setup workdir.
ENV WORKDIR /app
WORKDIR $WORKDIR
