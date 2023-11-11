#!/bin/bash
sudo docker build -t cassiozareck/little-ecom-auth:latest .
sudo docker push cassiozareck/little-ecom-auth:latest

kubectl rollout restart deploy auth