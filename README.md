![Go](https://img.shields.io/badge/Go-1.21-blue)
![Docker](https://img.shields.io/badge/Docker-Enabled-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supported-blue)
![Build](https://img.shields.io/badge/Build-Passing-brightgreen)
![License](https://img.shields.io/badge/License-MIT-green)

# 🏬 microstore 
microstore is a modular e-commerce backend portfolio project built with Go using a microservices architecture. It demonstrates clean architecture principles, inter-service communication via REST, and containerized deployments using Docker and Docker Compose.

## 🔧 Features

- Auth Service: User registration, login, JWT generation, basic admin logic.
- Product Service: CRUD operations for products with stock tracking.
- API Gateway: Central routing with request forwarding, CORS config, and middleware for authentication and admin access.
- PostgreSQL for persistent storage, easily swappable with in-memory for tests.
- Unit-tested services with in-memory repos for isolated logic testing.
- Ratelimiter for limiting request per windows time base on the ip
- Caching product for demonstration of cahsing mechanisim.
- Searching product base on title or description.
- seprate shared library for seprating the microservice logic from common and repetetive ...
- ready made postman collection for interacting with API.
- Dockerized using docker-compose for local orchestration.

## Structure Overview
```graphql
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
## Architecture

- Language: GO (Golang)
- Design: Clean Architecture + Modular Packages
- API Communication: REST (HTTP)
- Databases: PostgreSQL
- Caching: Redis
- Search: Zincsearch (product indexing)
- Containerizaion: Docker + Docker Compose

## Getting Started
### Prerequisites
- Go
- Docker
- Make

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

## Service Details

### Auth Service (auth-service)
- User signup/login
- JWT authenticaiton
- Password hashing

### Product Service (product-service)
- add/list/delete/search products
- add product with product validator
- in-memory or PostgreSQL repository support
- zincsearch-based indexing
- redis-base product cache

### Order Service (order-service)
- create and fetch orders 
- create order validator

### API Gateway (gateway)
- Rotues public HTTP traffic
- Simple routing using GO and chi
- Handles requests for auth, product, and order endpoint
- Redus-base Ratelimiter

## License

This project is for educational and portfolio purposes. Feel free to use it as a reference or learning resource.

