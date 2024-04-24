#start from golang base image
FROM golang:1.22-alpine

#install dependencies
RUN apk update
RUN apk add --no-cache git build-base

#setup workdir
RUN mkdir -p /app
ENV WORKDIR /app
WORKDIR $WORKDIR

COPY ./src/user-service/ ./src/user-service/
COPY .env ./
COPY go.mod ./
COPY go.sum ./
COPY Makefile ./

# install dependencies project
WORKDIR $WORKDIR/src/user-service/
RUN go mod download

# run the Go app.
WORKDIR $WORKDIR
CMD ["make","start-user"]