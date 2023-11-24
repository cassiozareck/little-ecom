# Little E-Commerce Notifier Service

## Overview

The Little E-Commerce Notifier is a microservice designed to send notifications based on events in an e-commerce system. It listens for messages on specific RabbitMQ queues and processes them to send out email notifications to users.

## RabbitMQ Integration

The notifier service integrates with RabbitMQ, a robust and scalable message broker, to receive messages about various events such as item additions and purchases. It establishes a connection to the RabbitMQ server and declares queues for different notification types.

### Queues

- `ecom-queue-item-added`: This queue receives messages when a new item is added to the e-commerce platform. The notifier service listens to this queue and sends an email to the user, informing them about the new item.
- `ecom-queue-item-bought`: This queue receives messages when an item is purchased. The notifier service listens to this queue and sends a confirmation email to the buyer.

## Kubernetes Deployment

The notifier service is containerized using Docker and can be deployed on a Kubernetes cluster. The provided `notifier.yaml` file describes a Kubernetes Deployment that manages the lifecycle of the notifier pods. The deployment ensures that the notifier service is always running and can scale as needed to handle the load.

The `deploy.sh` script is used to build the Docker image, push it to a container registry, and trigger a rolling update of the notifier deployment on Kubernetes. This ensures zero downtime deployments and seamless updates to the notifier service.

## Conclusion

The Little E-Commerce Notifier is a critical component of the e-commerce system, responsible for ensuring that users are kept informed about important events through email notifications. Its integration with RabbitMQ for message handling and Kubernetes for deployment makes it a resilient and scalable solution for real-time notification delivery.
