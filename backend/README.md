## Backend API Overview

The backend API of the application is responsible for handling CRUD operations related to items (products) and storing them in a database. It is designed to run as a containerized service within a Kubernetes (k8s) cluster, and it interacts with other services such as authentication (auth) and notification (notifier) services.

### Interaction with Auth Service

The backend API uses JWT (JSON Web Tokens) for securing endpoints. It interacts with the auth service to validate tokens sent by clients. When a request is made to a protected endpoint, the token is extracted from the `Authorization` header and sent to the auth service for validation. If the token is valid, the request is processed; otherwise, an error is returned.

### Interaction with Notifier Service

The notifier service is used to send notifications (such as emails) to users. The backend API publishes messages to a RabbitMQ exchange with a specific routing key when certain actions are performed, such as adding or buying an item. The notifier service consumes these messages and sends notifications accordingly.

### Middlewares

The backend API uses two middlewares:

- **Logging Middleware**: Logs all incoming HTTP requests with their method and URL path.
- **CORS Middleware**: Sets the necessary CORS headers to allow cross-origin requests, which is essential for the frontend to interact with the API from a different domain.

### Endpoints

The API provides several endpoints for item management:

- `GET /item/{id}`: Retrieve a specific item by its ID.
- `POST /item`: Add a new item.
- `DELETE /item/{id}`: Remove an item by its ID.
- `PUT /item/{id}`: Update an item by its ID.
- `GET /items/{owner}`: Retrieve all items owned by a specific user.
- `GET /items`: Retrieve all items.
- `POST /buy/{id}`: Buy an item by its ID.

### Kubernetes Manifest

The `backend.yaml` file is a Kubernetes manifest that defines the deployment and service for the backend API. It specifies the container image to use, the number of replicas, and the service type for external access. The deployment ensures that the desired state of the backend API is maintained within the k8s cluster.

### How It Works

The backend API is built as a Docker image and pushed to Docker Hub. The `deploy.sh` script automates this process. The Kubernetes cluster pulls the image from Docker Hub and runs it as specified in the `backend.yaml` manifest. The API connects to MongoDB for data persistence and RabbitMQ for messaging with the notifier service.

The backend is designed to be scalable and resilient, running multiple replicas for load balancing and high availability. It is also stateless, which makes it suitable for running in a containerized environment like Kubernetes.

### Scalability and Resilience

The backend is designed to be scalable and resilient. It can handle a growing number of requests by scaling the number of pods in the Kubernetes cluster. The stateless nature of the services allows for easy scaling without the need for persistent local storage.

