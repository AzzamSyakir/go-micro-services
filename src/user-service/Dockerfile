#start from golang base image
FROM golang:1.22-alpine

#install dependencies
RUN apk update
RUN apk add --no-cache git build-base

#setup workdir
RUN mkdir -p /app
ENV WORKDIR /app
WORKDIR $WORKDIR

#copy source from host to workdir inside container
COPY . .
COPY .env .env
COPY go.mod  go.mod
COPY go.sum  go.sum

# Build the Go app.
RUN go build -a -o User  src/user-service/cmd/web/main.go
#run go app
CMD ["./User"]