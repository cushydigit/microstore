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
- Rate limiter to limit requests per time window based on IP address
- Product caching to demonstrate caching mechanism
- Product search based on title or description
- Separate shared library to isolate reusable and common logic from microservices
- Ready-made Postman collection for interacting with the API
- Dockerized using docker-compose for local orchestration.

## 📁 Structure Overview
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
## 🧩 Architecture

- Language: GO (Golang)
- Design: Clean Architecture + Modular Packages
- API Communication: REST (HTTP)
- Databases: PostgreSQL
- Caching: Redis
- Search: Zincsearch (product indexing)
- Containerizaion: Docker + Docker Compose

## 🚀 Getting Started
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

## 📚 Service Details

### 🔐 Auth Service (`auth-service`)
- User signup/login
- JWT authenticaiton
- Password hashing

### 🛒 Product Service (`product-service`)
- Add/list/delete/search products
- Add product with product validator
- In-memory or PostgreSQL repository support
- Zincsearch-based indexing
- Redis-based product cache

### 📦 Order Service (`order-service`)
- Create and fetch orders 
- Create order validator

### 🌐 API Gateway (`gateway`)
- Rotues public HTTP traffic
- Simple routing using GO and chi
- Handles requests for auth, product, and order endpoint
- Redus-base Ratelimiter

## 📬 Postman Collection

You can use the Postman collection to test all available endpoints:
[Download Postman Collection](./postman/microstore.api..postman_collection.json)
[![Run in Postman](https://run.pstmn.io/button.svg)](https://www.postman.com/material-astronaut-37601285/cushydigit/folder/w8ksi5h/microstore-api?action=share&creator=21076955&ctx=documentatio)

## 📜 License

This project is for educational and portfolio purposes. Feel free to use it as a reference or learning resource.

