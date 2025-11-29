Articles API Helm Chart

This Helm chart deploys the Articles REST API (Go) along with a PostgreSQL database and a full monitoring stack
(Prometheus & Grafana) on Kubernetes.


🏗 Architecture

The chart deploys the following components:

*   **Articles API**: A Go application handling REST requests.
*   **PostgreSQL**: Managed via the Bitnami dependency chart.
*   **Prometheus**: Metrics collection (via kube-prometheus-stack).
*   **Grafana**: Visualization dashboards (via kube-prometheus-stack).


🚀 Quick Start

**1. Add Dependency Repositories**

This chart relies on standard community charts for the database and monitoring.

```sh
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
```

**2. Download Dependencies**

Fetch the specific versions of Postgres and Prometheus defined in `Chart.yaml`.

```sh
helm dependency update .
```

**3. Install the Chart**

Install the release with the name `articles-release`.

```sh
helm upgrade --install articles-release . -f values.yaml
```

⚙️ Configuration

You can override values via the `--set` flag or by creating your own `my-values.yaml`.

### Common Parameters

| Parameter | Description | Default |
| :--- | :--- | :--- |
| `image.repository` | Application image | `antohachaban/articles-go` |
| `image.tag` | Application version | `0.0.1` |
| `replicaCount` | Number of API replicas | `1` |
| `appEnv` | Environment (e.g., production, dev) | `production` |
| `service.port` | Service port | `8080` |
| `resources` | CPU/Memory limits | `Requests: 100m/128Mi` |

### Database & Secrets

| Parameter | Description | Default |
| :--- | :--- | :--- |
| `secrets.postgresPassword` | Password for DB User | `supersecretpassword` |
| `postgresql.auth.database` | Database name to create | `articles` |
| `postgresql.persistence.size` | PVC Size | `1Gi` |

**Example: Custom Installation**

To run with a specific password and 3 replicas:

```sh
helm upgrade --install articles-release . \
  --set replicaCount=3 \
  --set secrets.postgresPassword="my-hard-password" \
  --set postgresql.auth.password="my-hard-password"
```

🖥️ Accessing the Application

By default, the application runs as `ClusterIP` (internal only). To access it from your local machine:

Port Forward the service:

```sh
kubectl port-forward svc/articles-release-app 8080:8080
```

Test the API:

```sh
# Create Article
curl -X POST http://localhost:8080/article \
  -H "Content-Type: application/json" \
  -d '{"title": "Hello Kubernetes"}'

# Health Check
curl http://localhost:8080/health
```