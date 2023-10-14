#!/bin/bash
# Script to construct entire cluster. This assumes you have: helm, minikube and also that minikube isn't running and
# its empty

# Start Minikube
echo "Starting Minikube..."
minikube start --cpus 2 --memory 4096
if [ $? -ne 0 ]; then
    echo "Error starting Minikube. Exiting."
    exit 1
fi

# Update and add Helm repositories
echo "Updating Helm repositories..."
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Setup mongo DB
if ! helm list --deployed | grep -q "^mongo"; then
    echo "Installing MongoDB..."
    helm install mongo bitnami/mongodb -f mongo-helm.yaml
else
    echo "MongoDB release already exists. Skipping installation."
fi

# Fundamental component
echo "Setting up backend..."
kubectl apply -f backend.yaml

# Setup EFK stack
echo "Setting up EFK stack..."
kubectl apply -f elastic-search.yaml
kubectl apply -f fluentd.yaml
kubectl apply -f kibana.yaml

# Install RabbitMQ
if ! helm list --deployed | grep -q "^rabbitmq"; then
    echo "Installing RabbitMQ..."
    helm install rabbitmq bitnami/rabbitmq
else
    echo "RabbitMQ release already exists. Skipping installation."
fi

echo "Setup complete."
