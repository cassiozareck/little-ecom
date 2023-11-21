### Little-ecom

### How it works

### Visual intuition

### Requirements for execution:
- Minikube installed
- Helm installed
- Python and pygame (optional, for the cluster visualizer)

### How to execute
Run `./setup_minikube.sh` on your shell, and it will automatically build the cluster.

### How to interact with it it
1- Open: https://editor.swagger.io/
2- Copy and paste openapi.yaml content into it

This will let you interact with the cluster, the request will reach first at ingress which will redirect
to the backend service. You can send any request listed

3- (optional) run: `python3 cluster_visualizer.py`. This will open a cluster renderer, which blinks circles representing
pods through new logs and request flow. It can be useful to see if one request correctly traverses k8s cluster flow 
(e.g., INGRESS -> Backend -> Mongo).

## Utils
Under /scripts you can run the debugger script (which will open rabbitmq panel, mongodb compass, 
and minikube dashboard). Under the same folder, you can also run the start_efk.sh script which will start 
an efk stack, so you can get cluster logs in real-time using Grafana for visualization, Elasticsearch 
for storage, filtering and search, and Fluentd as an agent to get logs in every node.

## Common commands

- `helm repo add bitnami https://charts.bitnami.com/bitnami` -> Adds bitnami repo
- `helm install mongo bitnami/mongodb -f helms-manifests/mongo-helm.yaml` -> Installs mongodb
- `minikube dashboard` -> Opens the k8s dashboard
- `k apply -f helms-manifests/rabbitmq.yaml` -> Applies rabbitmq manifest
- `k rollout restart deployment/backend` -> Restarts backend deployment
