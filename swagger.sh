
# Get Minikube IP
MINIKUBE_ADDRESS=$(minikube ip)

# Replace placeholder with actual IP in OpenAPI template
sed "s/{{MINIKUBE_IP}}/http:\/\/$MINIKUBE_ADDRESS/g" openapi_template.yaml > openapi.yaml

swagger serve -F=swagger openapi.yaml
