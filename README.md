# Local install

```bash
docker build -t hello-world-server:latest .
```

Update values.yaml with:
```yaml
image:
  repository: hello-world-server
  tag: latest
  pullPolicy: Never
```

```bash
helm install hello-world ./helm
```

# Port Forward
Forward the service to your local machine
```bash
kubectl port-forward service/hello-world-hello-world-server 8080:80
```

# Then open browser to:
# http://localhost:8080

# For k3s-specific considerations:
If you want to use k3s's built-in load balancer (Traefik), you can enable ingress:

# Update helm/values.yaml:
# ingress.enabled: true
# ingress.hosts[0].host: hello-world.your-domain.com

helm upgrade hello-world ./helm
