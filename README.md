# PreBoot
A Go-based CLI tool that diagnoses local setup failures in Go repositories by running deterministic, YAML-defined checks.

## Development

```bash
make build       # build ./preboot
make test        # run tests with race detector
make ci          # build + vet + test (mirrors CI)
make help        # list all targets
```
