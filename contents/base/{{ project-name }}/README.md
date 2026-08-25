# {{ project_title }}

A basic Go service using the chi router. It ships with the platform's deployment,
observability, and CI/CD plumbing but no domain logic — grow it into whatever you need.

The service exposes a single stub endpoint on the service port:

```
GET /  ->  {"service":"{{ project-name }}","status":"ok"}
```

Health, readiness, liveness, and Prometheus metrics are served on the management port.

## Setup

```bash
make setup
```

## Development

```bash
go run ./cmd/server
```

## Build

```bash
make build
```

## Test

```bash
make test
```
