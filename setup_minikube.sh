#!/bin/bash

# Script to start Minikube and set up the base components

# Start Minikube
echo "Starting Minikube..."
minikube start --cpus 2 --memory 4096 --nodes=2

if [ $? -ne 0 ]; then
    echo "Error starting Minikube. Exiting."
    exit 1
fi

# Update and add Helm repositories
echo "Updating Helm repositories..."
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

# Install NGINX Ingress Controller
if ! helm list --deployed | grep -q "^nginx-ingress"; then
    echo "Installing NGINX Ingress Controller..."
    helm install nginx-ingress ingress-nginx/ingress-nginx -f ingress.yaml
else
    echo "NGINX Ingress Controller release already exists. Skipping installation."
fi

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

# Install RabbitMQ
echo "Setting up rabbitmq..."
kubectl apply -f rabbitmq.yaml

# Notifier
echo "Setting up notifier..."
kubectl apply -f notifier.yaml

echo "Base setup complete."
