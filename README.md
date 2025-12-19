# 🚀 Ultimate Full-Stack Template (Go + Next.js)

A production-ready, high-performance boilerplate designed for scalability, maintainability, and developer experience. Built with **Go (Fiber)** for the backend and **Next.js 14** for the frontend, following **Clean Architecture** principles.

## ✨ Features

### 🔐 Authentication & Security
*   **OTP-based Authentication**: Secure login flow using SMS (Senator provider).
*   **JWT Authorization**: Stateless authentication with access and refresh tokens.
*   **RBAC Ready**: Structure in place for Role-Based Access Control.
*   **Security Headers**: Helmet middleware, CORS configuration, and CSRF protection.
*   **Rate Limiting**: Redis-backed rate limiting to prevent abuse.

### 💳 Payments & Financials
*   **Multi-Gateway Support**:
    *   **Zarinpal**: Full implementation (Request & Verify).
    *   **Vandar**: Full implementation (Request & Verify).
    *   **Card-to-Card**: Manual receipt submission and tracking.
*   **Transaction Tracking**: Unified transaction model for all payment methods.

### 📱 Communication
*   **SMS Gateway**: Modular adapter pattern.
    *   **Senator**: Implemented provider for sending OTPs.
    *   *Easy to add more providers (e.g., KavehNegar, Twilio).*

### 📦 Storage & Data
*   **S3 Compatible Storage**: AWS SDK v2 integration for file uploads (MinIO, AWS S3, Cloudflare R2).
*   **PostgreSQL**: Primary relational database with `pgx` driver.
*   **Redis**: Caching, session management, and rate limiting.

### 🎨 Frontend Architecture (The "Perfect Stack")
*   **Next.js 14**: App Router, Server Components, and Server Actions.
*   **Shadcn UI**: Beautiful, accessible, and customizable components.
*   **State Management**:
    *   **TanStack Query**: Server state management, caching, and optimistic updates.
    *   **Nuqs**: URL-based state management for shareable filters and dialogs.
*   **Form Handling**: React Hook Form + Zod for robust validation.
*   **Visualizations**: Recharts for analytics dashboards.
*   **UX Polish**:
    *   **Sonner**: Stackable toast notifications.
    *   **Vaul**: Native-feeling mobile drawers.
    *   **Framer Motion**: Smooth animations.
    *   **Magic UI**: Landing page visual effects.

### 🛠 DevOps & Observability
*   **Dockerized**: Full stack containerization with Docker Compose.
*   **OpenTelemetry**: Distributed tracing instrumented across the backend.
*   **Logging**: Structured JSON logging with Zap (ready for Loki/ELK).
*   **CI/CD**: GitHub Actions pipeline for testing, linting, and building images.
*   **IaC**: Terraform scripts for provisioning AWS infrastructure (S3, RDS, ElastiCache).

---

## 🏗 Architecture

### Backend: Clean Architecture (Hexagonal)
The backend is structured to separate concerns and make the application testable and independent of frameworks.

```
backend/
├── cmd/api/            # Application entry point
├── internal/
│   ├── core/           # Pure business logic (Domain)
│   │   ├── domain/     # Entities and business rules
│   │   ├── ports/      # Interfaces (Input/Output ports)
│   │   └── services/   # Use cases implementation
│   ├── adapter/        # Implementation of ports (Infrastructure)
│   │   ├── handler/    # HTTP Handlers (Fiber)
│   │   ├── storage/    # Database repositories (Postgres, S3)
│   │   ├── payment/    # Payment gateway adapters
│   │   └── sms/        # SMS gateway adapters
└── pkg/                # Shared utilities (Logger, Telemetry)
```

### Frontend: Modular Feature-Based
```
frontend/
├── src/
│   ├── app/            # Next.js App Router pages
│   ├── components/     # UI Components
│   │   ├── ui/         # Shadcn base components
│   │   ├── dashboard/  # Dashboard specific widgets
│   │   └── ...
│   ├── lib/            # Utilities and helpers
│   ├── hooks/          # Custom React hooks
│   └── services/       # API clients
```

---

## 🚀 Getting Started

### Prerequisites
*   **Docker & Docker Compose** (Recommended)
*   **Go 1.21+** (For local backend dev)
*   **Node.js 20+** (For local frontend dev)

### Quick Start (Docker)

1.  **Clone the repository:**
    ```bash
    git clone <repo-url>
    cd full-template
    ```

2.  **Environment Setup:**
    ```bash
    cp .env.example .env
    # Edit .env with your credentials (DB, Redis, SMS, Payments)
    ```

3.  **Run with Docker Compose:**
    ```bash
    make docker-up
    # OR
    docker-compose up -d --build
    ```

4.  **Access the Application:**
    *   **Frontend**: `http://localhost:80` (via Nginx)
    *   **Backend API**: `http://localhost:80/api`
    *   **Swagger Docs**: `http://localhost:80/api/swagger/index.html`
    *   **Grafana**: `http://localhost:3000`
    *   **Jaeger**: `http://localhost:16686`

### Local Development

#### Backend
```bash
cd backend
go mod download
go run cmd/api/main.go
```

#### Frontend
```bash
cd frontend
npm install
npm run dev
```

---

## 🧪 Testing

### Backend
*   **Unit Tests**: `go test ./internal/core/...`
*   **Integration Tests**: `go test ./internal/test/integration/...` (Requires Docker)

### Frontend
*   **Unit Tests**: `npm test`
*   **E2E Tests**: `npx playwright test`

---

## 📦 Deployment

### CI/CD
A GitHub Actions workflow (`.github/workflows/ci.yml`) is configured to:
1.  Lint Backend (GolangCI-Lint) & Frontend (ESLint).
2.  Run Unit Tests.
3.  Build Docker Images.

### Infrastructure as Code
Terraform scripts in `terraform/` can be used to provision production infrastructure on AWS.

```bash
cd terraform
terraform init
terraform apply
```

---

## 📝 License
This project is licensed under the MIT License.
