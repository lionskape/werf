#!/bin/bash -e

script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
project_dir=$script_dir/../..

if [[ -z "$WERF_TEST_K8S_BASE64_KUBECONFIG" ]]; then
  echo "script requires \$WERF_TEST_K8S_BASE64_KUBECONFIG" >&2
  exit 1
fi

if [[ -z "$WERF_TEST_K8S_DOCKER_REGISTRY" ]] || [[ -z "$WERF_TEST_K8S_DOCKER_REGISTRY_USERNAME" ]] || [[ -z "$WERF_TEST_K8S_DOCKER_REGISTRY_PASSWORD" ]]; then
  echo "script requires \$WERF_TEST_K8S_DOCKER_REGISTRY, \$WERF_TEST_K8S_DOCKER_REGISTRY_USERNAME and \$WERF_TEST_K8S_DOCKER_REGISTRY_PASSWORD" >&2
  exit 1
fi

echo $WERF_TEST_K8S_DOCKER_REGISTRY_PASSWORD | docker login $WERF_TEST_K8S_DOCKER_REGISTRY -u $WERF_TEST_K8S_DOCKER_REGISTRY_USERNAME --password-stdin

export KUBECONFIG=$(mktemp -d)/config
if [[ "$OSTYPE" == "linux-gnu" ]]; then
  echo $WERF_TEST_K8S_BASE64_KUBECONFIG | base64 -d -w0 > $KUBECONFIG
elif [[ "$OSTYPE" == "darwin"* ]]; then
  echo $WERF_TEST_K8S_BASE64_KUBECONFIG | base64 -D > $KUBECONFIG
fi
