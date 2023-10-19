
echo "Opening Minikube Dashboard..."
minikube dashboard &

echo "Exposing RabbitMQ Manager on a NodePort..."
kubectl expose service rabbitmq --type=NodePort --name=rabbitmq-manager --port=15672
RABBITMQ_PORT=$(kubectl get svc rabbitmq-manager -o=jsonpath='{.spec.ports[0].nodePort}')
echo "Opening RabbitMQ Manager..."
xdg-open "http://$(minikube ip):$RABBITMQ_PORT/" &

# Port-forward MongoDB from mongo-mongodb-0
echo "Port-forwarding MongoDB to localhost:27017..."
kubectl port-forward pod/mongo-mongodb-0 27017:27017 &
sleep 5  # Give it a few seconds to establish the port-forward

MONGO_URL="mongodb://localhost:27017/"

echo "Opening MongoDB Compass with the connection string..."
mongodb-compass $MONGO_URL &

