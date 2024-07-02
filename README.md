# Go Microservices 

POS-API develop using Microservices architecture pattern 
# Microservices

## What are Microservices?
```Microservices``` are an architectural approach to developing an application as a collection of small, self-contained services that each fulfill a specific business purpose. Developers can build these services in several programming languages, deploy, scale, and maintain them independently, and enable communication between them via well-defined APIs. The following image demonstrates how ```Microservices``` work in practice

![Microservices architecture](https://firebasestorage.googleapis.com/v0/b/image-to-onlin.appspot.com/o/diagram.png?alt=media&token=bcbba1dc-b579-4b71-ae08-359802b19a34)

As shown in the image above, clients (mobile, web, or desktop applications) send requests to an API gateway, which serves as the entry point, routing each request to the appropriate Microservices. Furthermore, each service operates independently, interacting with its own database and, if necessary, with other Microservices or an external API or service to fulfill requests.

Microservices vs monoliths
To further understand microservices, it's helpful to contrast them with the traditional pattern of developing applications — the monolithic architecture. Applications in a monolithic architecture are often constructed in layers, e.g., a presentation layer to handle user interaction, a business logic layer to process data according to business rules, and a data access layer to communicate with the database.


## introduction
This project uses the ``Microservices`` architecture as explained above, this ``project aims to be a boilerplate/template and can be used as an example and guide`` to make it easier for developers to create micro-services applications with Golang
## Features
- docker container per services
- db per services 
- communication between services with grpc
- user services 
- product services 
- order services
- auth services 

## Installation

### first clone this repo
```bash
git clone "github repo links"
```
### build and run  projects 
```bash
go mod tidy
```
```bash
make start-docker
```
### run test
```bash
make start-test
```
