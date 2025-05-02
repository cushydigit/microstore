![Go](https://img.shields.io/badge/Go-1.21-blue)
![Docker](https://img.shields.io/badge/Docker-Enabled-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supported-blue)
![Build](https://img.shields.io/badge/Build-Passing-brightgreen)
![License](https://img.shields.io/badge/License-MIT-green)
# 🏬 microstore 
Microstore is a modular, production-style e-commerce backend built with Go using a clean microservices architecture. It demonstrates real-world patterns for scalable systems with service separation, JWT-based authentication, and centralized API routing.

## 🔧 Features

- Auth Service: User registration, login, JWT generation, basic admin logic.

- Product Service: CRUD operations for products with stock tracking.

- API Gateway: Central routing with request forwarding, CORS config, and middleware for authentication and admin access.

- PostgreSQL for persistent storage, easily swappable with in-memory for tests.

- Unit-tested services with in-memory repos for isolated logic testing.

- Dockerized using docker-compose for local orchestration.

## 🧪 Goals
- Showcase clean Go code, modular design, and testability.

- Avoid unnecessary frontend — everything is tested using Postman or similar tools.

- Serve as a professional backend portfolio project to demonstrate understanding of Go, microservices, API design, and system integration.
