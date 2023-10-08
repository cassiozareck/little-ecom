#!/bin/bash
sudo docker build -t cassiozareck/hello-kube:latest .
sudo docker push cassiozareck/hello-kube:latest

kubectl rollout restart deployment hello-kube

