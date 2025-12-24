# Relearn Golang - refresh my knowledge of golang

## 🚀 Key Features

*   **Clean Architecture**: Separation of concerns into Handler, Service, Domain, and Repository layers.
*   **High Performance**: Built with [Gin Web Framework](https://github.com/gin-gonic/gin) for speed.
*   **Secure Authentication**: JWT-based authentication with Access and Refresh token rotation.
*   **Database**: PostgreSQL integration using [GORM](https://gorm.io/).
*   **Containerized**: Fully Dockerized for consistent development and deployment environments.
*   **Developer Friendly**:
    *   **Hot Reloading**: Integrated with [Air](https://github.com/air-verse/air) for instant feedback.
    *   **Swagger Documentation**: Auto-generated API docs using [swaggo](https://github.com/swaggo/swag).
    *   **Make Utility**: comprehensive `Makefile` for common tasks.

## 🛠 Tech Stack

*   **Language**: Go (Golang) 1.21+
*   **Framework**: Gin
*   **Database**: PostgreSQL
*   **ORM**: GORM
*   **Auth**: JWT (JSON Web Tokens)
*   **Docs**: Swagger (OpenAPI 2.0)
*   **Tooling**: Air, Docker, Docker Compose

## 📦 Getting Started

### Prerequisites

*   Go 1.21 or higher
*   Docker & Docker Compose
*   Make (optional, but recommended)

### Quick Start (Local Development)

1.  **Clone the repository**
    ```bash
    git clone https://github.com/prachaya-orr/relearn-golang.git
    cd relearn-golang
    ```

2.  **Start Infrastructure (PostgreSQL)**
    ```bash
    make docker-up
    ```

3.  **Run the Application**
    ```bash
    make run
    # OR for hot-reloading
    make watch
    ```

    The API will be available at `http://localhost:8080`.

4.  **View Documentation**
    Access the Swagger UI at:
    `http://localhost:8080/swagger/index.html`

## 🧪 Testing

Run quality assurance tests with:

```bash
make test
```

## 🏗 Project Structure

```text
.
├── cmd/                # Entry points (api, migration tools)
├── internal/
│   ├── domain/         # Business entities and interfaces
│   ├── handler/        # HTTP Handlers (Controllers)
│   ├── middleware/     # Gin Middlewares (Auth, Logging)
│   ├── repository/     # Data Access Layer
│   └── service/        # Business Logic Layer
├── docs/               # Swagger generated docs
└── ...
```

## 📝 API Endpoints Summary

*   **Auth**:
    *   `POST /signup`: Register a new user
    *   `POST /login`: Authenticate and get tokens
    *   `POST /refresh-token`: Rotate access tokens

*(See Swagger docs for full list)*
g