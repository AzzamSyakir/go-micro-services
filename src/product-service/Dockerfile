#start from golang base image
FROM golang:1.25-alpine

#install dependencies
RUN apk update
RUN apk add --no-cache git build-base

#setup workdir
RUN mkdir -p /app
ENV WORKDIR=/app
WORKDIR $WORKDIR

RUN go install github.com/air-verse/air@latest


COPY ./src/product-service/ ./src/product-service/
COPY grpc/ grpc/
COPY product.air.toml .
COPY .env .
COPY go.mod .
COPY go.sum .
COPY Makefile .
COPY tmp .

# install dependencies project
WORKDIR $WORKDIR/src/product-service/
RUN go mod download

# run the Go app.
WORKDIR $WORKDIR
CMD ["air", "-c", "product.air.toml"]
