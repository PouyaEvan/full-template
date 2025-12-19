# Full-Template Helm Chart

This Helm chart deploys the Full-Stack Template (Go + Next.js) application to Kubernetes.

## Prerequisites

- Kubernetes 1.21+
- Helm 3.0+
- PV provisioner support in the underlying infrastructure (for PostgreSQL and Redis persistence)

## Installing the Chart

To install the chart with the release name `my-release`:

```bash
# Add Bitnami repository for dependencies
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Install the chart
helm install my-release ./helm/full-template
```

## Configuration

See [values.yaml](values.yaml) for the full list of configurable parameters.

### Common Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `global.environment` | Environment name | `production` |
| `backend.enabled` | Enable backend deployment | `true` |
| `backend.replicaCount` | Number of backend replicas | `2` |
| `frontend.enabled` | Enable frontend deployment | `true` |
| `frontend.replicaCount` | Number of frontend replicas | `2` |
| `postgresql.enabled` | Deploy PostgreSQL | `true` |
| `redis.enabled` | Deploy Redis | `true` |

### Secrets

Update the following secrets in production:

```yaml
secrets:
  databaseUrl: "postgresql://user:pass@host:5432/db"
  jwtSecret: "your-secure-secret"
  zarinpalMerchantId: "your-merchant-id"
  vandarApiKey: "your-api-key"
```

### Ingress

Configure ingress for external access:

```yaml
ingress:
  enabled: true
  hosts:
    - host: your-domain.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: your-tls-secret
      hosts:
        - your-domain.com
```

## Upgrading

```bash
helm upgrade my-release ./helm/full-template
```

## Uninstalling

```bash
helm uninstall my-release
```

## Architecture

The chart deploys the following components:

- **Backend**: Go Fiber API server
- **Frontend**: Next.js application
- **PostgreSQL**: Primary database (optional, via Bitnami chart)
- **Redis**: Caching and session storage (optional, via Bitnami chart)
- **Ingress**: External access routing

```
┌─────────────────────────────────────────────────────────────┐
│                        Ingress                               │
└─────────────────────────────────────────────────────────────┘
                             │
              ┌──────────────┼──────────────┐
              │              │              │
              ▼              ▼              ▼
        ┌─────────┐    ┌─────────┐    ┌─────────┐
        │ Frontend│    │ Backend │    │   WS    │
        │ Service │    │ Service │    │ Service │
        └─────────┘    └─────────┘    └─────────┘
              │              │
              ▼              ▼
        ┌─────────┐    ┌─────────┐
        │PostgreSQL│    │  Redis  │
        └─────────┘    └─────────┘
```

## License

MIT
