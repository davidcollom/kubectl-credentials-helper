# kubectl-credentials-keychain

This is a `kubectl` credentials helper. It is using your local OS keychain to store sensitive content of your `KUBECONFIG`.

To install:

```bash
go install github.com/davidcollom/kubectl-credentials-keychain@latest
```

You can also download this from releases.

To secure your existing `KUBECONFIG`, use:

```bash
kubectl-credentials-keychain secure --kubeconfig path
```

If `--kubeconfig` was not provided - it will try to find `KUBECONFIG` env variable and as a last resort - default user home `~/.kube/config`.

This will save all sensitive info from a local `KUBECONFIG` to your OS specific keychain. Sensitive considered the following:

- `.user.username`
- `.user.password`
- `.user.client-certificate-data`
- `.user.client-key-data`

For user entries that it finds, it will replace in your `KUBECONFIG` with the following:

```yaml
apiVersion: v1
kind: Config
users:
- name: <name>
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1
      command: kubectl-credentials-helper
      provideClusterInfo: true
      interactiveMode: Never
```

Now, every time you are trying to access this cluster - the helper will fetch sensitive info from OS specific keychain.

## Development and Testing

### Prerequisites

- Go 1.24.5 or later
- Make

### Building

```bash
# Build the binary
make build

# Install dependencies
make deps

# Clean build artifacts
make clean
```

### Testing

The project has been refactored to be more testable with proper dependency injection and interfaces:

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage
make test-coverage

# Run tests for specific packages
go test ./cmd ./internal/logger -v
```

### Code Quality

```bash
# Format code
make fmt

# Run go vet
make vet

# Run linter (requires golangci-lint)
make lint

# Run all checks
make check
```

### Logging

The project uses structured logging with configurable levels:

- Set `KUBECTL_CREDENTIALS_HELPER_DEBUG=true` for debug output
- Logs are written to stderr to keep stdout clean for kubectl integration
- Uses logrus for structured logging with timestamps

### Architecture

The code has been restructured for better testability:

- `cmd/credential_helper.go`: Core business logic with dependency injection
- `cmd/root.go`: CLI command setup
- `internal/logger/`: Structured logging package with interfaces
- Interfaces allow for easy mocking in tests
- Separated concerns for better unit testing
