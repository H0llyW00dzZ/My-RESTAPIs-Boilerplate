# K8s Deployment for REST API Boilerplate - MySQL

<p align="center">
   <img src="https://kubernetes.io/images/kubernetes.png" alt="fully-managed-and-isolated-by-k8s" width="80">
</p>

This repository contains Kubernetes (K8s) deployment files for a MySQL database service. The deployment includes configurations for MySQL, and Vertical Pod Autoscaler (VPA).

## Prerequisites

Before deploying the MySQL service, ensure that you have the following:

- A running Kubernetes cluster
- `kubectl` command-line tool installed and configured to communicate with your cluster
- A container image for MySQL (e.g., `mysql:8.0`)

## Deployment

To deploy the MySQL service using the provided K8s deployment files, follow these steps:

1. Update the `mysql-deployment.yaml` file with your desired configuration, such as the MySQL image, resource limits, and environment variables.

2. Create the necessary secrets for the deployment:

   ```bash
   kubectl create secret generic mysql-root-pass --from-literal=password=your-mysql-root-password
   kubectl create secret generic mysql-ssl --from-file=certificate.cer=path/to/your/certificate.cer --from-file=server-key.pem=path/to/your/server-key.pem --from-file=ECC.crt=path/to/your/ECC.crt
   ```

   Replace `your-mysql-root-password` with your desired MySQL root password, and `path/to/your/certificate.cer`, `path/to/your/server-key.pem`, and `path/to/your/ECC.crt` with the paths to your SSL certificate files.

3. Create a persistent volume claim for MySQL storage:

   ```bash
   kubectl apply -f mysql-storage.yaml
   ```

   Update the `mysql-storage.yaml` file with your desired storage configuration (e.g., size, storage class).

4. Apply the deployment file to your Kubernetes cluster:

   ```bash
   kubectl apply -f mysql-deployment.yaml
   ```

   This command will create the necessary namespace, ConfigMap, deployment, service, and VPA for the MySQL service.

5. Wait for the deployment to complete and the pods to be in the "Running" state:

   ```bash
   kubectl get pods -n database
   ```

6. Access the MySQL service using the external IP or domain name assigned to the service:

   ```bash
   kubectl get service mysql-service -n database
   ```

   The output will display the external IP or domain name for the MySQL service.

## Monitoring and Scaling

The provided deployment includes a Vertical Pod Autoscaler (VPA) configuration that automatically adjusts the resource limits (CPU and memory) for the MySQL pod based on usage. You can monitor the VPA and the deployment using the following commands:

```bash
kubectl get vpa -n database
kubectl get deployment -n database
```

Adjust the VPA configuration in the `mysql-deploy.yaml` file to suit your application's scaling requirements.

## Customization

The provided deployment files are designed to be customizable. You can modify the resource limits, environment variables, and other configurations according to your application's needs.

## Cleanup

To remove the deployed resources from your Kubernetes cluster, run the following commands:

```bash
kubectl delete -f mysql-deploy.yaml
```

This will delete the deployment, service, VPA, and persistent volume claim associated with the MySQL service.
