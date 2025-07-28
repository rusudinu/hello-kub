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

Deployment workflow:
Option 1: Copy just the helm folder

# From your dev machine
scp -r helm/ user@k3s-server:/path/to/deployment/

# On k3s server
helm install hello-world ./helm/ --kubeconfig /etc/rancher/k3s/k3s.yaml

Option 2: Clone repo but only use helm folder

# On k3s server
git clone your-repo.git
cd your-repo
helm install hello-world ./helm/
Option 3: Package the helm chart



# From your dev machine
helm package ./helm/
# This creates hello-world-server-0.1.0.tgz

# Copy and install the package
scp hello-world-server-0.1.0.tgz user@k3s-server:
helm install hello-world hello-world-server-0.1.0.tgz
