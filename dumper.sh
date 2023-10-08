#!/bin/bash

# Directory to save the logs
LOG_DIR="./logs"
mkdir -p "$LOG_DIR"

while true; do
    # List all deployments in the current namespace
    for deployment in $(kubectl get deployments -o=jsonpath='{.items[*].metadata.name}'); do
        # For each deployment, get the pods
        for pod in $(kubectl get pods -l app="$deployment" -o=jsonpath='{.items[*].metadata.name}'); do
            # Dump the logs to a file named <deployment>-<pod>.log
            kubectl logs "$pod" > "$LOG_DIR/$deployment-$pod.log" 2>&1 &
        done
    done
    # Wait for a certain time interval before fetching the logs again
    sleep 60
done
