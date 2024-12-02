# K8s Deployment for REST API Boilerplate - Redis Insight

<p align="center">
   <img src="https://kubernetes.io/images/kubernetes.png" alt="sailing-with-k8s" width="80">
   <img src="https://i.imgur.com/wGetVaj.png" alt="The-Black-Pearl" width="80">
</p>

This repository contains Kubernetes (K8s) deployment files for a REST API boilerplate application. The deployment includes a Redis Insight service and a Redis Insight deployment.

## Prerequisites

Before deploying the application, ensure that you have the following:

- A running Kubernetes cluster
- `kubectl` command-line tool installed and configured to communicate with your cluster
- Redis Insight Docker image (`redis/redisinsight:latest`) available in your container registry

## Deployment

To deploy the REST API boilerplate application using the provided K8s deployment files, follow these steps:

1. Update the `redis-insight.yaml` file with your desired configuration, such as the number of replicas, resource limits, and environment variables.

2. Create the necessary secrets for the deployment:

   ```bash
   kubectl create secret generic redisinsight-tls-secrets --from-file=tls.key=path/to/your/tls.key --from-file=tls.crt=path/to/your/tls.crt
   kubectl create secret generic redisinsight-encryption-secret --from-literal=encryption-key=your-encryption-key
   ```

   Replace `path/to/your/tls.key`, `path/to/your/tls.crt`, and `your-encryption-key` with your actual TLS key, TLS certificate, and encryption key, respectively.

3. Apply the deployment files to your Kubernetes cluster:

   ```bash
   kubectl apply -f redisinsight-deployment.yaml
   ```

   This command will create the Redis Insight service and deployment in your Kubernetes cluster.

4. Wait for the deployment to complete and the pods to be in the "Running" state:

   ```bash
   kubectl get pods -l app=redisinsight
   ```

5. Access the Redis Insight service using the external IP or domain name assigned to the service:

   ```bash
   kubectl get services redis-insight-service
   ```

   The output will display the external IP or domain
