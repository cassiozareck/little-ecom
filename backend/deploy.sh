#!/bin/bash
docker build -t cassiozareck/little-ecom-backend:latest .
docker push cassiozareck/little-ecom-backend:latest

kubectl rollout restart deploy backend