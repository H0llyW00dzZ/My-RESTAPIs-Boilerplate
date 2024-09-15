# K8s Deployment for REST API Boilerplate - Smooth Sailing ⛵ ☸

<p align="center">
   <img src="https://kubernetes.io/images/kubernetes.png" alt="fully-managed-and-isolated-by-k8s" width="80">
</p>

This repository contains Kubernetes (K8s) deployment files for a REST API boilerplate application. The deployment includes configurations for the REST API service, Ingress controller, and Horizontal Pod Autoscaler (HPA).

## Prerequisites

Before deploying the application, ensure that you have the following:

- A running Kubernetes cluster
- `kubectl` command-line tool installed and configured to communicate with your cluster
- A container image for the REST API application

> [!NOTE]
> Since this deployment supports `100% HPA`, which is suitable for handling billions of requests/workers (combined with worker package) `concurrently and efficiently`, it is recommended not to attach it with any `storage (PVC, PV)` to this deployment.
> This is because HPA is not `100% compatible` if the deployment has `storage (PVC, PV) attached due to its limitations`, unless you build your own `storage mechanism` that can be shared among multiple pods (e.g., `capable of up to 1K Pods or more that consider supports 100% HPA`) while this deployment handles billions of requests/workers `concurrently and efficiently`.

## Deployment

To deploy the REST API boilerplate application using the provided K8s deployment files, follow these steps:

1. Update the `restapis-deploy.yaml` file with your desired configuration, such as the number of replicas, resource limits, and environment variables. Replace `<IMAGE_HERE>` with the actual container image for your REST API application.

2. Create the necessary secrets for the deployment:

   ```sh
    ./create_k8s_secret.sh
   ```

   Replace the placeholders with your actual values for database connections, Redis connections, timeouts, and other environment variables required by your REST API application.

3. Apply the deployment file to your Kubernetes cluster:

   ```bash
   kubectl apply -f restapis-deploy.yaml
   ```

   This command will create the necessary namespace, deployment, service, and HPA for the REST API application.

4. Update the `restapis-ingress.yaml` file with your desired configuration, such as the host name, TLS certificate, and paths for the REST API routes and frontend routes.

5. Apply the Ingress configuration to your Kubernetes cluster:

   ```bash
   kubectl apply -f restapis-ingress.yaml
   ```

   This command will create the Ingress resource for routing traffic to the REST API service.

6. Wait for the deployment to complete and the pods to be in the "Running" state:

   ```bash
   kubectl get pods -n restapis
   ```

7. Access the REST API service using the configured host name and paths specified in the Ingress configuration.

## Monitoring and Scaling

The provided deployment includes a Horizontal Pod Autoscaler (HPA) configuration that automatically scales the number of replicas based on CPU and memory utilization. You can monitor the HPA and the deployment using the following commands:

```bash
kubectl get hpa -n restapis
kubectl get deployment -n restapis
```

Adjust the HPA configuration in the `restapis-deploy.yaml` file to suit your application's scaling requirements.

## Customization

The provided deployment files are designed to be customizable. You can modify the resource limits, environment variables, and other configurations according to your application's needs. Additionally, you can adjust the Ingress configuration to match your desired routing rules and TLS settings.

## Cleanup

To remove the deployed resources from your Kubernetes cluster, run the following commands:

```bash
kubectl delete -f restapis-ingress.yaml
kubectl delete -f restapis-deploy.yaml
```

This will delete the Ingress resource, deployment, service, and HPA associated with the REST API application.
