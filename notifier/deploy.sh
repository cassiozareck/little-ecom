#!/bin/bash
sudo docker build -t cassiozareck/little-ecom-notifier:latest .
sudo docker push cassiozareck/little-ecom-notifier:latest

kubectl rollout restart deployment notifier-deployment

