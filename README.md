
# Go Micro-Services 

template boilerplate for go project using micro-services architecture 
# micro-services
## What are microservices?
```Microservices``` are an architectural approach to developing an application as a collection of small, self-contained services that each fulfill a specific business purpose. Developers can build these services in several programming languages, deploy, scale, and maintain them independently, and enable communication between them via well-defined APIs. The following image demonstrates how ```microservices``` work in practice.

![Micro-Service architecture](https://firebasestorage.googleapis.com/v0/b/image-to-onlin.appspot.com/o/micro-services-architecture.png?alt=media&token=346717a2-f2cc-4aaa-a73c-d1bcd73b9de2)
As shown in the image above, clients (mobile, web, or desktop applications) send requests to an API gateway, which serves as the entry point, routing each request to the appropriate microservice. Furthermore, each service operates independently, interacting with its own database and, if necessary, with other microservices or an external API or service to fulfill requests.

Microservices vs monoliths
To further understand microservices, it's helpful to contrast them with the traditional pattern of developing applications — the monolithic architecture. Applications in a monolithic architecture are often constructed in layers, e.g., a presentation layer to handle user interaction, a business logic layer to process data according to business rules, and a data access layer to communicate with the database.


## introduction
This project uses the ``micro-services`` architecture as explained above, this ``project aims to be a boilerplate/template and can be used as an example and guide`` to make it easier for developers to create micro-services applications with Golang
## Features

- db per services 
- communication between services
- user services 
- product services 
- order services
- api gateway services 


## Installation

### first clone this repo
```bash
git clone "githu repo links"
```
### build and run docker compose 
```bash
docker compose -f ./docker/docker-compose.yml build --no-cache
```
```bash
docker compose -f ./docker/development/docker-compose.yml up -d
```
### run project
```bash
go mod tidy
```
```bash
go run cmd/web/main.go
```
