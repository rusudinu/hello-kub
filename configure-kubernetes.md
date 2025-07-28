# Kubernetes Configuration Prompt

Create a complete Kubernetes configuration for a Go web application with the following specifications:

## Variables to Replace
- **APP_NAME**: The application name (e.g., hello-kub)
- **IMAGE**: Container image (e.g., ghcr.io/rusudinu/hello-kub:latest)
- **CONTAINER_PORT**: Application port inside container (e.g., 8080)
- **DOMAIN**: Your domain name (e.g., hello-kub.rusudinu.ro)
- **NAMESPACE**: Kubernetes namespace (e.g., default)

## Required Resources

### 1. Deployment
- Name: APP_NAME
- 2 replicas initially
- Container image: IMAGE
- Container port: CONTAINER_PORT
- Resource limits: 100m CPU, 128Mi memory
- Resource requests: 50m CPU, 64Mi memory
- Image pull policy: Always
- Health checks: liveness and readiness probes on HTTP GET `/` targeting CONTAINER_PORT
- Labels: app.kubernetes.io/name=APP_NAME, app.kubernetes.io/instance=APP_NAME, app.kubernetes.io/version=1.0.0

### 2. Service
- Name: APP_NAME
- Type: ClusterIP
- Port 80 targeting container port CONTAINER_PORT (named "http")
- Selector matches deployment labels

### 3. Horizontal Pod Autoscaler (HPA)
- Name: APP_NAME
- Min replicas: 2, Max replicas: 10
- Scale based on CPU utilization: 40% average
- Scale based on memory utilization: 40% average
- Target the deployment: APP_NAME

### 4. HTTPS Ingress with Redirect
- Name: APP_NAME
- Ingress class: traefik
- Host: DOMAIN
- Path: / (Prefix match)
- TLS enabled with Let's Encrypt certificate (secret: APP_NAME-tls)
- Annotations:
  - cert-manager.io/cluster-issuer: letsencrypt-prod
  - traefik.ingress.kubernetes.io/router.entrypoints: web,websecure
  - traefik.ingress.kubernetes.io/router.middlewares: NAMESPACE-APP_NAME-redirect-https@kubernetescrd

### 5. Traefik Middleware
- Name: APP_NAME-redirect-https
- Type: redirectScheme with scheme=https, permanent=true
- This forces HTTP traffic to redirect to HTTPS
- **Important**: Reference in ingress as `NAMESPACE-APP_NAME-redirect-https@kubernetescrd` (namespace-name format)

## File Structure
Create separate YAML files:
- deployment.yaml
- service.yaml  
- hpa.yaml
- ingress.yaml
- middleware.yaml

## Requirements
- All resources should have consistent labeling using APP_NAME
- Remove Helm-specific labels (helm.sh/chart, app.kubernetes.io/managed-by)
- Use standard Kubernetes labels only
- Ensure proper resource relationships and selectors match
- Replace all variables (APP_NAME, IMAGE, CONTAINER_PORT, DOMAIN, NAMESPACE) with actual values