#!/bin/bash
docker build -t cassiozareck/little-ecom-notifier:latest .
docker push cassiozareck/little-ecom-notifier:latest

kubectl rollout restart deploy notifier-deployment
