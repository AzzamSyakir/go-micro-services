#start from golang base image
FROM golang:1.22-alpine

#install dependencies
RUN apk update
RUN apk add --no-cache git build-base

#setup workdir
RUN mkdir -p /app
ENV WORKDIR /app
WORKDIR $WORKDIR

COPY . .
# install dependencies project
RUN go mod download

# run the Go app.
CMD [ "make", "start-product" ] 