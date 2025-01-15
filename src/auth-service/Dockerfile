#start from golang base image
FROM golang:1.23-alpine

#install dependencies
RUN apk update
RUN apk add --no-cache git build-base

#setup workdir
RUN mkdir -p /app
ENV WORKDIR /app
WORKDIR $WORKDIR

RUN go install github.com/air-verse/air@latest

# copy source from host to docker container
COPY ./src/auth-service/ ./src/auth-service/
COPY grpc/ grpc/
COPY auth.air.toml .
COPY .env .
COPY go.mod .
COPY go.sum .
COPY Makefile .

# install dependencies project
WORKDIR $WORKDIR/src/auth-service/
RUN go mod download

# run the Go app.
WORKDIR $WORKDIR
CMD ["air", "-c", "auth.air.toml"]
