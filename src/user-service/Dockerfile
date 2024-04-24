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
CMD [ "make", "start-user" ] 

# nanti coba yng ini
# FROM golang:1.22-alpine

# WORKDIR /app

# RUN go install github.com/cosmtrek/air@latest
# COPY . .
# COPY .env .env
# COPY go.mod  go.mod
# COPY go.sum  go.sum

# RUN go mod download

# CMD ["air", "-c", ".air.toml"] 