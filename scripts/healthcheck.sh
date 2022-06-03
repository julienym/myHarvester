#!/bin/bash
echo -n $KUBE64 | base64 -d > /tmp/.cluster
export KUBECONFIG="/tmp/.cluster"
while true; do
  kubectl get ns 2>/dev/null
  if [[ $? -ne 0 ]]
  then
    echo "Waiting for API to be ready..."
    sleep 15s
  else
    exit 0
  fi
done