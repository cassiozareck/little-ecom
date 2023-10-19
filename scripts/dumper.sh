#!/bin/bash

LOG_DIR="./k8s_logs"
mkdir -p $LOG_DIR

while true; do
  # Fetch logs from all pods in the default namespace
  for POD in $(kubectl get pods -o jsonpath='{.items[*].metadata.name}'); do
    # Fetch logs with verbosity level 7 and limit to the last 500 lines
    kubectl logs $POD --v=7 | tail -n 500 > "${LOG_DIR}/${POD}.log"
  done

  # Wait for 20 seconds before the next iteration
  sleep 20
done
