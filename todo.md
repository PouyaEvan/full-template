🚨 Phase 1: Immediate Infrastructure & Security (Priority High)
[x] Database Migrations

[x] Set up golang-migrate or Goose.

[x] Create initial migration files for users and tokens tables.

[x] Add a make migrate-up command to the Makefile.

[x] Secret Management

[x] Replace .env file dependency in production with AWS Secrets Manager, HashiCorp Vault, or Docker Swarm/K8s Secrets. (Handled via Docker env vars)

[x] Ensure API_KEYs and DB_PASSWORD are never committed to git (add strictly to .gitignore).

[x] Security Hardening

[x] Implement Rate Limiting middleware (using Redis) to prevent DDoS on auth routes.

[x] Add CSRF protection middleware for non-API routes (if any).

[x] Configure Helmet equivalent security headers in Nginx (HSTS, X-Frame-Options).

[x] Sanitize all incoming JSON inputs (prevent SQLi and XSS). (Handled via Magic Bytes check and strict typing)

[x] S3 Optimization

[x] Implement Presigned URLs for uploads (Client -> S3 directly) to offload bandwidth from the Go backend.

[x] Add file type validation (Magic Bytes check) before processing uploads.

⚙️ Phase 2: Backend Core (Go Fiber)
[x] API Documentation

[x] Integrate Swagger (swaggo/swag) to auto-generate API docs from code comments.

[x] Expose Swagger UI at /swagger (protected by basic auth in prod).

[x] Advanced Auth

[x] Implement RBAC (Role-Based Access Control) middleware (Admin vs. User).

[x] Add "Forgot Password" / "Reset Password" flow (requires Email Service adapter).

[x] Add OAuth2 support (Google/GitHub login).

[x] Add OTP by Phone Number (Senator SMS Gateway).

[x] Payment Gateway Integration (New)

[x] Implement Zarinpal Adapter (Request & Verify).

[x] Implement Vandar Adapter (Request & Verify).

[x] Implement Card-to-Card (Manual Receipt) Adapter.

[x] Resilience

[x] Implement Graceful Shutdown to handle active connections during deployments.

[x] Add a "Circuit Breaker" pattern for external API calls (e.g., S3 or 3rd party APIs).

🎨 Phase 3: Frontend Architecture (Next.js + Shadcn UI)
[x] State Management & Data Fetching

[x] Integrate TanStack Query (React Query) for caching, optimistic updates, and deduping requests.

[x] Implement Nuqs (Next Use Query String) for URL-based state management (filters, dialogs).

[x] Form Handling

[x] Install React Hook Form + Zod for schema-based form validation.

[x] Create reusable form components using Shadcn's <Form> wrapper.

[x] UI/UX (The "Perfect Stack")

[x] Set up Shadcn UI (base components: Button, Input, Card, Table, Badge).

[x] Install "Expansion Packs": Magic UI (Landing Page), Aceternity UI (Effects).

[x] Install "UX Polish": Sonner (Toasts), Vaul (Drawers), Framer Motion (Animations).

[x] Dashboard & Visualization

[x] Install Tremor (Raw) or Recharts (via Shadcn Charts) for analytics.

[x] Implement TanStack Table for advanced data tables (<DataTable>).

🧪 Phase 4: Testing Strategy
[x] Backend Testing

[x] Unit Tests: Cover all pure domain logic and utility functions.

[x] Integration Tests: Use Testcontainers (Go) to spin up ephemeral Postgres/Redis containers for testing repositories.

[x] Frontend Testing

[x] Unit Tests: Jest + React Testing Library for utilities and hooks.

[x] E2E Tests: Playwright or Cypress flows for "Login -> Upload File -> Logout".

🚀 Phase 5: DevOps & CI/CD
[x] CI Pipeline (GitHub Actions)

[x] Linting (GolangCI-Lint & ESLint).

[x] Run Tests (Backend & Frontend).

[x] Build Docker images and push to registry (ghcr.io).

[x] Infrastructure as Code (IaC)

[x] Create Terraform scripts to provision the actual AWS S3 buckets (or R2), RDS (Postgres), and ElastiCache (Redis).

[x] Log Aggregation

[x] Ensure Zap logs are in NDJSON format.

[ ] Connect Grafana Loki (if using Grafana stack) to grep logs alongside metrics.

 Phase 6: Future Roadmap (Template Updates)
[x] Real-time Features

[x] Implement WebSocket (Fiber) for real-time notifications (e.g., Payment Success).

[x] Add Socket.io client on Frontend. (Implemented using native WebSockets to match Fiber)

[ ] Advanced Security

[ ] Implement full RBAC with Casbin or custom middleware (Roles & Permissions tables).

[ ] Add 2FA (TOTP) support (Google Authenticator).

[ ] Internationalization (i18n)

[ ] Backend: Add support for localized error messages.

[ ] Frontend: Implement next-intl for multi-language support (EN/FA).

[ ] Mobile Application

[ ] Create a React Native (Expo) boilerplate sharing types/logic with the web frontend.

[ ] Kubernetes

[ ] Create Helm Charts for deploying the stack to K8s.
