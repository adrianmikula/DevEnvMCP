# Integration Tests

This directory contains integration tests that use Testcontainers to test the Dev-Env Sentinel with real build tools (Maven, npm) running in Docker containers.

## Prerequisites

1. **Docker**: Docker must be installed and running on your machine
2. **Go 1.24+**: Required for building and running tests
3. **Testcontainers**: Automatically installed via `go.mod`

## Running Integration Tests

### Run All Integration Tests

```bash
go test -tags=integration ./integration -v
```

### Run Specific Test

```bash
go test -tags=integration ./integration -v -run TestIntegration_RealMavenBuild
```

### Skip Integration Tests

Integration tests are skipped when:
- The `-short` flag is used: `go test -short ./...`
- Docker is not available
- Test data is missing

## Test Structure

### Container-Based Tests

These tests spin up Docker containers with real build tools:

- **TestIntegration_DetectMavenProject**: Tests Maven container setup
- **TestIntegration_DetectNpmProject**: Tests npm container setup  
- **TestIntegration_RealMavenBuild**: Tests actual Maven compilation
- **TestIntegration_RealNpmBuild**: Tests npm project setup
- **TestIntegration_BuildFreshnessWithRealMavenProject**: Tests build freshness detection with real builds

### Local File System Tests

These tests use local testdata without containers:

- **TestIntegration_EcosystemDetectionWithRealProject**: Tests ecosystem detection with local testdata

## Test Data

Integration tests use:
- `testdata/` directory for local file system tests
- Docker containers for isolated build tool testing

## CI/CD Integration

In CI environments, ensure:
1. Docker is available (or tests will be skipped)
2. Sufficient resources for container execution
3. Network access for pulling Docker images

Example GitHub Actions setup:

```yaml
- name: Run integration tests
  run: go test -tags=integration ./integration -v
  env:
    DOCKER_HOST: unix:///var/run/docker.sock
```

## Troubleshooting

### Tests Timeout

If tests timeout, increase the container startup timeout in test setup functions.

### Docker Not Available

Tests will skip automatically if Docker is not available. To verify Docker:

```bash
docker ps
```

### Container Cleanup

Containers are automatically cleaned up with `AutoRemove: true`. If containers persist, manually remove them:

```bash
docker ps -a | grep testcontainers
docker rm -f <container-id>
```

## Performance

Integration tests are slower than unit tests:
- Container startup: ~5-10 seconds per test
- Build execution: ~2-5 seconds per build
- Total: ~30-60 seconds for full suite

Use `-short` flag to skip integration tests during fast feedback loops.

