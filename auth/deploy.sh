#!/bin/bash
docker build -t cassiozareck/little-ecom-auth:latest .
docker push cassiozareck/little-ecom-auth:latest

kubectl rollout restart deploy auth