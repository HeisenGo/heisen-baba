# Project Name

## Overview

**HeseinBaba** is an online hotel and tour reservation system built with Go that leverages modern microservices architecture. It integrates several technologies to provide a scalable, maintainable, and high-performance system. This project uses Traefik for reverse proxy, Consul for service discovery, RabbitMQ for messaging, and various other technologies to ensure seamless communication and data management.

## Architecture

### Microservices

The project is built using a microservices architecture, with services following a combination of Clean Architecture and Hexagonal Architecture principles. This ensures a clear separation of concerns and maintains flexibility and scalability.

### Key Components

- **Traefik**: Acts as a reverse proxy, managing routing and load balancing between services.
- **Consul**: Handles service registry and discovery.
- **gRPC**: Used for efficient communication between services.
- **RabbitMQ**: Facilitates asynchronous messaging and communication between services.
- **Fiber**: Framework used to build REST APIs.
- **JWT**: Handles authentication and authorization.
- **Redis**: Provides caching to enhance performance.
- **PostgreSQL**: Main relational database used for data storage.

## Getting Started

### Prerequisites

Ensure you have the following installed:
- Docker (for containerization)
- Docker Compose (for managing multi-container Docker applications)

### Setup

1. **Clone the repository:**

    ```bash
    git clone https://github.com/HeisenGo/heisen-baba.git
    ```

2. **Build and run services using Docker Compose:**

    ```bash
    make dockerize
    ```

    This command will build all necessary Docker images and start the services defined in `docker-compose.yml`.


### Configuration

All configuration is managed through environment variables. You can customize your setup by editing the `config.yaml` file or setting environment variables directly Or rename `config.yaml.example` to `config.yaml` in each service.
