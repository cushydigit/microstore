<h1 align="center">
    🏬 microstore 
</h1>
<p align="center">
  <a href="https://github.com/cushydigit/microstore/LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-green.svg" alt="License MIT">
  </a>
  <img src="https://img.shields.io/badge/Go-1.24-blue.svg" alt="React 18+">
  <img src="https://img.shields.io/badge/Build-Passing-brightgreen.svg" alt="Project Status">
  <img src="https://img.shields.io/badge/Docker-Enabled-blue.svg" alt="Project Status">
  <img src="https://img.shields.io/badge/Build-Passing-brightgreen.svg" alt="Project Status">
  <img src="https://img.shields.io/badge/PostgreSQL-Supported-blue.svg" alt="Project Status">
</p>
<p align="center">
microstore is a modular e-commerce backend portfolio project built with Go using a microservices architecture. It demonstrates clean architecture principles, inter-service communication via REST, and containerized deployments using Docker and Docker Compose.
</p>

## Features

- Auth Service: User registration, login, JWT auth, basic admin logic
- Product Service: CRUD, stock tracking, search, and caching
- Order Service: Order creation and retrieval with validation
- API Gateway: Centralized routing with CORS, authentication, and admin middleware
- PostgreSQL (with optional in-memory mode) for persistence
- Redis Caching for faster product access
- ZincSearch-based product search by title/description
- Rate Limiting by IP to control request flow
- Shared Library for reusable helpers, middleware, DB, cache, and utils
- Postman Collection included for easy API testing
- Dockerized with docker-compose for local orchestration


## Technology Stack 

### Core Technologies

- **Language**: Go 1.24
- **Architecture**: Clean Architecture with microservices
- **Communication**: REST APIs over HTTP
- **Containerizaion**: Docker + Docker Compose

### Data Storage

- **Primary Database**: PostgreSQL with seperate schemas per service
- **Caching Layer**: Redis for product cashing and rate limiting 
- **Seach Engine**: Zincsearch for product full-text search

### Developement Tools

- **Build System**: Makefile with common commands
- **API Tesing**: Postman collection included
- **Testing**: In-memory repositories for unit-testing

## Getting Started 🚀
The project requires `Go` ,`Docker` and `Make` to get started.

clone the repo
```bash
git clone https://github.com/cushydigit/microstore.git
cd microstore

```

run all services
```bash
make up

```

stop services
```bash
make down

```

## Service Responsibilites

### Gateway Service

The `gateway` service serves as the single entry point, handling:

- Request routing to appropriate business sevices
- JWT authentication validation
- CORS configuration
- IP-based rate limiting via Redis

### Business Sevices

Each business service owns its domain and data:

- `auth-service`: User management and JWT token generation
- `product-service`: Product lifecycle with advanced search and caching
- `order-service`: Order processing with product validation

### Shared Components

The `shared` module provides common functionality:

- Type definitions for inter-service communication
- Database connection management
- Redis client abstraction
- Zincsearch client abstraction
- Middleware for authenticaion and validation

## Structure Overview

```tree
microstore/
├── auth-service         # Handles user registration, login, authentication, and rate limitiing
├── product-service      # Manages products, supports search and caching
├── order-service        # Manages customer orders and order workflows
├── gateway              # API Gateway routing external HTTP traffic to services
├── shared               # Common utilities, middleware, DB, Redis, search clients
├── db                   # SQL migrations for initializing databases
├── docker-compose.yml   # Orchestrates all services with PostgreSQL & Redis
├── Makefile             # Common build and run commands

```

## Request Processing Flow

![App Screenshot](./assets/microstore_flow_overview.png)

## Postman Collection

You can use the Postman collection to test all available endpoints:
[Download Postman Collection](./postman/microstore.api.postman_collection.json)
[![Run in Postman](https://run.pstmn.io/button.svg)](https://www.postman.com/material-astronaut-37601285/cushydigit/folder/w8ksi5h/microstore-api?action=share&creator=21076955&ctx=documentatio)

## License

This project is for educational and portfolio purposes. Feel free to use it as a reference or learning resource.

