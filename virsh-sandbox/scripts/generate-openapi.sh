#!/usr/bin/env bash
set -euo pipefail

swag init --dir .,./internal/ansible,./internal/diff,./internal/error,./internal/rest,./internal/vm,./internal/workflow --generalInfo ./cmd/api/main.go --parseDependency --parseInternal
