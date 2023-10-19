#!/bin/bash
# Script to construct entire cluster. This assumes you have: helm, minikube and also that minikube isn't running and
# its empty

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
helm repo add elastic https://helm.elastic.co
helm repo update

# Install NGINX Ingress Controller
if ! helm list --deployed | grep -q "^nginx-ingress"; then
    echo "Installing NGINX Ingress Controller..."
    helm install nginx-ingress ingress-nginx/ingress-nginx
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

# Setup Elasticsearch with Helm
if ! helm list --deployed | grep -q "^elasticsearch"; then
    echo "Installing Elasticsearch..."
    helm install elastic bitnami/elasticsearch -f elastic-helm.yaml
else
    echo "Elasticsearch release already exists. Skipping installation."
fi

# Setup Fluentd with Helm
if ! helm list --deployed | grep -q "^fluentd"; then
    echo "Installing Fluentd..."
    helm install fluentd oci://registry-1.docker.io/bitnamicharts/fluentd -f fluentd-helm.yaml
else
    echo "Fluentd release already exists. Skipping installation."
fi

# Setup Kibana with Helm
if ! helm list --deployed | grep -q "^kibana"; then
    echo "Installing Kibana..."
    helm install kibana elastic/kibana -f kibana-helm.yaml
else
    echo "Kibana release already exists. Skipping installation."
fi

# Install RabbitMQ
if ! helm list --deployed | grep -q "^rabbitmq"; then
    echo "Installing RabbitMQ..."
    helm install rabbitmq bitnami/rabbitmq -f rabbitmq-helm.yaml
else
    echo "RabbitMQ release already exists. Skipping installation."
fi

# Notifier
echo "Setting up notifier..."
kubectl apply -f notifier.yaml

echo "Setup complete."
