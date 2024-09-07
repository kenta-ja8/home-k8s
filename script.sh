#!/bin/bash

set -euxo pipefail

function buildAndPush() {
  docker buildx build \
    --platform linux/arm64 \
    --provenance false \
    --tag ghcr.io/kenta-ja8/$1:latest \
    --push \
    --build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
    --build-arg TARGET=$2 \
    ./src/
}

if [ "$#" -eq 0 ]; then
  echo "Usage: $0 <function_name>"
  exit 1
fi

case "$1" in
  buildServer)
      echo "buildServer"
      buildAndPush home-k8s-app-server ./cmd/server
      ;;
  buildJob)
      echo "buildJob"
      buildAndPush home-k8s-app-job ./cmd/job
      ;;
  buildAll)
      echo "buildAll"
      buildAndPush home-k8s-app-server ./cmd/server
      buildAndPush home-k8s-app-job ./cmd/job
      ;;
  *)
      echo "Unknown function: $1"
      exit 1
      ;;
esac
