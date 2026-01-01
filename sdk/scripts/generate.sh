#!/bin/bash
# scripts/generate.sh

set -e

echo "Generating SDK..."

# Merge OpenAPI specs first
npx openapi-merge-cli --config openapi/openapi-merge.json

# Generate with custom templates
docker run --rm \
  -v ${PWD}:/local \
  openapitools/openapi-generator-cli generate --skip-validate-spec \
  -i /local/openapi/combined.yaml \
  -g python \
  -o /local/virsh-sandbox-py/ \
  -c /local/.openapi-generator/config.yaml \
  -t /local/.openapi-generator/templates/python/

echo "Running polish script..."
python3 scripts/polish_sdk.py

echo "Running tests..."
cd virsh-sandbox-py
pip install -r requirements.txt
black .
isort .
mypy virsh_sandbox
pip install -r test-requirements.txt
pytest

echo "Finished!"
