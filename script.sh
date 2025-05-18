#!/bin/bash

set -euxo pipefail

cleanup() {
  echo "Cleaning up..."
  kill $(jobs -p)
  wait
}
trap cleanup SIGINT

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
  portForward)
      echo "portForward"
      PORT_FORWARDS=(
        "kubectl port-forward svc/argocd-server -n argocd --address 0.0.0.0 8080:443"
        "kubectl port-forward svc/postgres-svc 15432:5432 --namespace app"
        "kubectl port-forward svc/grafana-svc 13000:3000 --namespace visualization"
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
