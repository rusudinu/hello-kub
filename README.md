# Deployment

## Local Development

Build and load image into kind:
```bash
docker build -t hello-kub:latest .
kind load docker-image hello-kub:latest
```

Update `kubernetes/deployment.yaml` image to use local build:
```yaml
image: "hello-kub:latest"
imagePullPolicy: Never
```

Deploy to local cluster:
```bash
kubectl apply -f kubernetes/
```

## Production Deployment

Deploy all resources:
```bash
kubectl apply -f kubernetes/
```

## Common Operations

### Port Forward
Forward the service to your local machine:
```bash
kubectl port-forward service/hello-kub 8080:80
```
Then open browser to: http://localhost:8080

### Update Deployment
After making changes to manifests:
```bash
kubectl apply -f kubernetes/
```

### Scale Application
```bash
kubectl scale deployment hello-kub --replicas=5
```

### View Resources
```bash
# Check deployment status
kubectl get deployments
kubectl get pods
kubectl get services
kubectl get ingress

# View logs
kubectl logs -f deployment/hello-kub

# Describe resources for troubleshooting
kubectl describe deployment hello-kub
kubectl describe pod <pod-name>
```

### Delete Resources
```bash
kubectl delete -f kubernetes/
```

### Rolling Updates
Update image tag in `kubernetes/deployment.yaml` then:
```bash
kubectl apply -f kubernetes/deployment.yaml
kubectl rollout status deployment/hello-kub
```

### Rollback
```bash
kubectl rollout undo deployment/hello-kub
```
