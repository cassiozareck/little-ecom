#!/bin/bash

kubectl set image deployment/backend little-ecom-backend=cassiozareck/little-ecom-backend:latest

kubectl rollout restart deployment backend



