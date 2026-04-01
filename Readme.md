# ShopCore

A production-ready e-commerce backend system built with Go, designed for high performance and maintainability.

## Overview

ShopCore provides a robust foundation for e-commerce platforms with comprehensive support for customer management, product catalogs, and order processing. Built with **Hexagonal Architecture** principles, it ensures clean separation between business logic and infrastructure concerns.

## Tech Stack

| Component | Technology |
|-----------|------------|
| **Language** | Go 1.25.6+ |
| **Web Framework** | Fiber v3 |
| **Database** | MongoDB |
| **Logging** | Zerolog |
| **Validation** | go-playground/validator v10 |
| **Serialization** | MessagePack |
| **Containerization** | Docker & Docker Compose |

## Architecture

ShopCore follows **Hexagonal Architecture** (Ports and Adapters) to achieve loose coupling and high cohesion:

```
┌─────────────────────────────────────────────────────────────────┐
│                        ADAPTERS                                  │
│  ┌─────────────────────┐          ┌─────────────────────────┐  │
│  │   Inbound           │          │   Outbound              │  │
│  │   (Driving)         │          │   (Driven)              │  │
│  │                     │          │                         │  │
│  │  • HTTP Handlers    │          │  • MongoDB Repository   │  │
│  │  • gRPC Services    │          │  • Email Service        │  │
│  │                     │          │  • Payment Gateway      │  │
│  └──────────┬──────────┘          └────────────┬────────────┘  │
│             │                                    │              │
└─────────────┼────────────────────────────────────┼──────────────┘
              │                                    │
              ▼                                    ▼
┌─────────────────────────────────────────────────────────────────┐
│                         PORTS                                    │
│  ┌─────────────────────┐          ┌─────────────────────────┐  │
│  │   Inbound Ports     │          │   Outbound Ports        │  │
│  │   (Interfaces)      │          │   (Interfaces)           │  │
│  │                     │          │                         │  │
│  │  • CustomerService  │          │  • CustomerRepository   │  │
│  │  • ProductService   │          │  • ProductRepository    │  │
│  │  • OrderService     │          │  • OrderRepository      │  │
│  └─────────────────────┘          └─────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────────────────────┐
│                         CORE                                     │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │   Domain Layer                                          │   │
│  │   • Entities (Customer, Product, Order, RunNumber)      │   │
│  │   • Business Rules & Logic                              │   │
│  └─────────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │   Service Layer                                         │   │
│  │   • Orchestration & Use Cases                            │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

## Features

### Customer Management
- Full CRUD operations
- Profile management with soft delete

### Product Catalog
- Product CRUD with stock tracking
- Automatic stock deduction on order

### Order Processing
- Complete order lifecycle management
- Order status workflow: `PENDING` → `PAID` → `SHIPPED` → `COMPLETED`
- Automatic order number generation with configurable prefix
- Product snapshot at order time for price consistency

### Auto Numbering System
- Configurable run number service for document IDs
- Format: `{PREFIX}-{DATE}-{SEQUENCE}`
- Atomic increment to prevent duplicate numbers

## Project Structure

```
shopcore/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
│
├── internal/
│   ├── adapters/
│   │   ├── inbound/
│   │   │   └── http/            # HTTP handlers & routes
│   │   │       ├── customer_handler.go
│   │   │       ├── product_handler.go
│   │   │       ├── order_handler.go
│   │   │       ├── run_number_handler.go
│   │   │       └── router.go
│   │   │
│   │   └── outbound/
│   │       └── mongodb/         # MongoDB repositories
│   │           ├── customer_repository.go
│   │           ├── product_repository.go
│   │           ├── order_repository.go
│   │           └── run_number_repository.go
│   │
│   ├── app/
│   │   └── app.go               # Application wiring
│   │
│   └── core/
│       ├── domain/              # Domain entities
│       │   ├── customer.go
│       │   ├── product.go
│       │   ├── order.go
│       │   └── run_number.go
│       │
│       ├── ports/
│       │   ├── inbound/        # Driving port interfaces
│       │   │   ├── customer_service.go
│       │   │   ├── product_service.go
│       │   │   ├── order_service.go
│       │   │   └── run_number_service.go
│       │   │
│       │   └── outbound/       # Driven port interfaces
│       │       ├── customer_repository.go
│       │       ├── product_repository.go
│       │       ├── order_repository.go
│       │       └── run_number_repository.go
│       │
│       └── services/           # Use case implementations
│           ├── customer_service.go
│           ├── product_service.go
│           ├── order_service.go
│           └── run_number_service.go
│
├── config/
│   ├── env.go                  # Environment configuration
│   ├── database.go             # MongoDB connection
│   └── framework.go            # Fiber setup
│
├── pkg/
│   └── helpers/                 # Shared utilities
│       ├── helper.go
│       └── time.go
│
├── scripts/                    # Helper scripts
├── assets/                     # Static assets
│
├── docker-compose.yaml         # Infrastructure setup
├── makefile                    # Build commands
├── go.mod                      # Go modules
└── go.sum                      # Dependencies lock
```

## API Endpoints

### Customers
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/customer` | Create customer |
| GET | `/api/v1/customer` | List customers |
| GET | `/api/v1/customer/:customer_id` | Get customer |
| PUT | `/api/v1/customer/:customer_id` | Update customer |
| DELETE | `/api/v1/customer/:customer_id` | Delete customer |

### Products
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/product` | Create product |
| GET | `/api/v1/product` | List products |
| GET | `/api/v1/product/:product_id` | Get product |
| PUT | `/api/v1/product/:product_id` | Update product |
| DELETE | `/api/v1/product/:product_id` | Delete product |

### Orders
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/order` | Create order |
| GET | `/api/v1/order` | List orders |
| GET | `/api/v1/order/:order_id` | Get order |
| PUT | `/api/v1/order/:order_id` | Update order |
| DELETE | `/api/v1/order/:order_id` | Cancel order |

### Run Numbers
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/run_number` | Initialize run number |

## Getting Started

### Prerequisites

- Go 1.25.6 or later
- Docker and Docker Compose (or Podman)

### Quick Start

1. **Start infrastructure services:**
   ```bash
   podman compose up -d
   ```

2. **Run the application:**
   ```bash
   make dev
   ```
   
   Or directly:
   ```bash
   go run cmd/api/main.go
   ```

3. **Access services:**
   - API: `http://localhost:3333`
   - MongoDB Express: `http://localhost:9040`

### Configuration

Environment variables (with defaults):
```bash
APP_PORT=3333
MONGO_DB_URI=mongodb://username:password@localhost:9041/
```

## Design Patterns

- **Hexagonal Architecture**: Clean separation between core business logic and external adapters
- **Repository Pattern**: Abstract data access layer
- **Dependency Injection**: Constructor-based DI for testability
- **Service Layer**: Orchestrates domain logic and repositories
- **DTO/Request-Response**: Clean API contracts

## Testing

```bash
go test ./...
```

## License

MIT
