#start from golang base image
FROM golang:1.22-alpine

#install dependencies
RUN apk update
RUN apk add --no-cache git build-base

#setup workdir
ENV WORKDIR /app
WORKDIR $WORKDIR

#copy source from host to workdir inside container
COPY ./src/user-service .
COPY .env .env
COPY go.mod go.mod
COPY go.sum go.sum

# Install the modules.
RUN go install -v ./...

# Clean all the downloaded modules.
RUN go clean -modcache

# Build the Go app.
RUN go build -C $WORKDIR/cmd/web -o $WORKDIR/build