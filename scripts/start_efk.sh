#!/bin/bash

# Script to deploy the EFK stack on an already running Minikube

helm repo add elastic https://helm.elastic.co
helm repo update
echo "Updating Helm repositories..."

# Setup Elasticsearch with Helm
if ! helm list --deployed | grep -q "^elasticsearch"; then
    echo "Installing Elasticsearch..."
    helm install elastic bitnami/elasticsearch -f helms/elastic-helm.yaml
else
    echo "Elasticsearch release already exists. Skipping installation."
fi

# Setup Fluentd with Helm
if ! helm list --deployed | grep -q "^fluentd"; then
    echo "Installing Fluentd..."
    helm install fluentd oci://registry-1.docker.io/bitnamicharts/fluentd -f helms/fluentd-helm.yaml
else
    echo "Fluentd release already exists. Skipping installation."
fi

# Setup Kibana with Helm
if ! helm list --deployed | grep -q "^kibana"; then
    echo "Installing Kibana..."
    helm install kibana elastic/kibana -f helms/kibana-helm.yaml
else
    echo "Kibana release already exists. Skipping installation."
fi

echo "EFK setup complete."
