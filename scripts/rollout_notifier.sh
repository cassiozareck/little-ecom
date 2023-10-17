kubectl set image deployment/notifier-deployment notifier-container=cassiozareck/little-ecom-notifier:latest

kubectl rollout restart deployment notifier-deployment