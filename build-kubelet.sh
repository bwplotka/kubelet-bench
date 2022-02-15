#!/bin/bash

set -e

K8S_ROOT=$1 # Path to local clone of https://github.com/kubernetes/kubernetes/

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Build kubelet.
cd ${K8S_ROOT}
make all WHAT=cmd/kubelet GOFLAGS=-v

cd ${DIR}

# Copy result.
cp ${K8S_ROOT}/_output/bin/kubelet .

# build local docker image out of it.
docker build -t kubelet:latest .