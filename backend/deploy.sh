#!/bin/bash
sudo docker build -t cassiozareck/little-ecom-backend:latest .
sudo docker push cassiozareck/little-ecom-backend:latest

kubectl rollout restart deployment little-ecom-backend

