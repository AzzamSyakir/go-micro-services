#start from golang base image
FROM golang:1.23-alpine

#install dependencies
RUN apk update
RUN apk add --no-cache git build-base
RUN go install github.com/air-verse/air@latest


#setup workdir
RUN mkdir -p /app
ENV WORKDIR /app
WORKDIR $WORKDIR

COPY ./src/order-service/ ./src/order-service/
COPY grpc/ grpc/
COPY order.air.toml .
COPY .env .
COPY go.mod .
COPY go.sum .
COPY Makefile .

# install dependencies project
WORKDIR $WORKDIR/src/order-service/
RUN go mod download

# run the Go app.
WORKDIR $WORKDIR
CMD ["air", "-c", "order.air.toml"]