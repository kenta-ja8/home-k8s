#!/bin/bash

set -euxo pipefail

if [ -z "${CONTAINER_CLI:-}" ]; then
  if command -v docker >/dev/null 2>&1; then
    CONTAINER_CLI=docker
  elif command -v podman >/dev/null 2>&1; then
    CONTAINER_CLI=podman
  else
    echo "Error: docker or podman is required" >&2
    exit 1
  fi
fi

cleanup() {
  echo "Cleaning up..."
  kill $(jobs -p)
  wait
}
trap cleanup SIGINT

function buildAndPush() {
  local image=$1
  local target=$2
  local cache_from_flags="${CACHE_FROM_FLAGS:-}"
  local cache_to_flags="${CACHE_TO_FLAGS:-}"
  local push_flag=""
  if [ -n "${cache_from_flags}" ]; then
    push_flag="--push"
  fi
  local build_date
  build_date=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

  if [ -n "${cache_from_flags}" ]; then
    cache_from_flags="${cache_from_flags},scope=${image}"
  fi

  if [ -n "${cache_to_flags}" ]; then
    cache_to_flags="${cache_to_flags},scope=${image}"
  fi

  "$CONTAINER_CLI" build \
    --platform linux/arm64 \
    --tag ghcr.io/kenta-ja8/"${image}":latest \
    --build-arg BUILD_DATE="${build_date}" \
    --build-arg TARGET="${target}" \
    ${push_flag} \
    ${cache_from_flags} \
    ${cache_to_flags} \
    --progress=plain \
    ./src/

  if [ -z "${push_flag}" ]; then
    # podmanの場合、--pushオプションがないためpushする
    "$CONTAINER_CLI" push ghcr.io/kenta-ja8/"${image}":latest
  fi
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
  portForward)
      echo "portForward"
      PORT_FORWARDS=(
        "kubectl port-forward svc/argocd-server -n argocd --address 0.0.0.0 8080:443"
        "kubectl port-forward svc/postgres-svc 15432:5432 --namespace app"
        "kubectl port-forward svc/grafana-svc 30090:3000 --namespace visualization"
      )
      for pf in "${PORT_FORWARDS[@]}"; do
        echo "Running: $pf"
        $pf &
      done
      wait
      ;;
  *)
      echo "Unknown function: $1"
      exit 1
      ;;
esac
