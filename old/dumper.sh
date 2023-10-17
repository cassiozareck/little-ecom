#!/bin/bash

# Directory to save the logs
LOG_DIR="./logs"
mkdir -p "$LOG_DIR"

# Verbosity level
VERBOSITY_LEVEL=5

while true; do
    # Get all running pods
    for pod in $(kubectl get pods --field-selector=status.phase=Running -o=jsonpath='{.items[*].metadata.name}'); do
        # Dump the logs to a file named <pod>.log
        kubectl logs "$pod" --v="$VERBOSITY_LEVEL" > "$LOG_DIR/$pod.log" 2>&1 &
    done
    # Wait for a certain time interval before fetching the logs again
    sleep 5
done
