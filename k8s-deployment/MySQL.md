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

1. Update the `mysql-deploy.yaml` file with your desired configuration, such as the MySQL image, resource limits, and environment variables.

2. Create the necessary secrets for the deployment:

   ```bash
   kubectl create secret generic mysql-root-pass --from-literal=password=your-mysql-root-password
   kubectl create secret generic mysql-ssl --from-file=certificate.cer=path/to/your/certificate.cer --from-file=server-key.pem=path/to/your/server-key.pem --from-file=ECC.crt=path/to/your/ECC.crt
   ```

   Replace `your-mysql-root-password` with your desired MySQL root password, and `path/to/your/certificate.cer`, `path/to/your/server-key.pem`, and `path/to/your/ECC.crt` with the paths to your SSL certificate files.

> [!NOTE]
> The example below demonstrates how `SSL/TLS` is correctly set up and can be utilized through any load balancer mechanism (recommended: a `standalone load balancer`). This deployment is designed to be dedicated and less noisy from neighboring services.
>
> Additionally, `TLS configuration` for `MySQL` is possible using `NGINX Ingress` Same `REST APIs Deployment` with other services (e.g., `TCP/UDP service NGINX`). However, this approach is not recommended. It is preferable to use a dedicated load balancer specifically for MySQL (e.g., `a single load balancer` only for MySQL) to avoid noisy from neighbor.
>
> It is advisable to use [`ECC`](https://en.wikipedia.org/wiki/Elliptic-curve_cryptography) instead of [`RSA`](https://en.wikipedia.org/wiki/RSA_(cryptosystem)) for encryption traffic; however, if `bandwidth is not a concern`, [`RSA`](https://en.wikipedia.org/wiki/RSA_(cryptosystem)) can be used as an alternative.


```bash
2024-09-13 15:52:16+00:00 [Note] [Entrypoint]: Entrypoint script for MySQL Server 8.4.2-1.el9 started.
2024-09-13 15:52:23+00:00 [Note] [Entrypoint]: Switching to dedicated user 'mysql'
2024-09-13 15:52:24+00:00 [Note] [Entrypoint]: Entrypoint script for MySQL Server 8.4.2-1.el9 started.
2024-09-13T15:52:27.782101Z 0 [System] [MY-015015] [Server] MySQL Server - start.
2024-09-13T15:52:29.412565Z 0 [System] [MY-010116] [Server] /usr/sbin/mysqld (mysqld 8.4.2) starting as process 1
2024-09-13T15:52:29.497857Z 1 [System] [MY-013576] [InnoDB] InnoDB initialization has started.
2024-09-13T15:52:33.304158Z 1 [System] [MY-013577] [InnoDB] InnoDB initialization has ended.
2024-09-13T15:52:33.543091Z 1 [System] [MY-011090] [Server] Data dictionary upgrading from version '80023' to '80300'.
2024-09-13T15:52:37.728885Z 1 [System] [MY-013413] [Server] Data dictionary upgrade from version '80023' to '80300' completed.
2024-09-13T15:52:52.425982Z 4 [System] [MY-013381] [Server] Server upgrade from '80039' to '80402' started.
2024-09-13T15:53:39.191347Z 4 [System] [MY-013381] [Server] Server upgrade from '80039' to '80402' completed.
2024-09-13T15:53:40.114222Z 0 [System] [MY-013602] [Server] Channel mysql_main configured to support TLS. Encrypted connections are now supported for this channel.
2024-09-13T15:53:40.288248Z 0 [System] [MY-010931] [Server] /usr/sbin/mysqld: ready for connections. Version: '8.4.2'  socket: '/var/lib/mysql/mysql.sock'  port: 3306  MySQL Community Server - GPL.
```
To connect to the MySQL server, you can use the domain name `database.example.com`.

#### Tested Connection on MySQL Workbench which works well and securely through a load balancer (as I am personally using a standalone load balancer that I made own):

<p align="center">
<img src="https://i.imgur.com/gHg4bLU.png" alt="tls-1-3-secure" width="500">
<img src="https://i.imgur.com/HKxXgjp.png" alt="tls-1-3-secure" width="500">
</p>

> [!NOTE]
> Also note that the default `MySQL configuration` from this repository uses `TLSv1.3` and can be used with `Public SSL/TLS CAs` as well. However, not all `MySQL clients` support the `TLSv1.3` protocol because they are built on `legacy systems`. For example, the `Navicat client's SSL/TLS protocol` for clients does not support `TLSv1.3`.

> [!TIP]
> When configuring TLS for `MySQL` and using a domain name through a load balancer, it is recommended to configure a network policy mechanism. The specifics of the network policy depend on the cloud provider. For instance, in GKE, when creating an external load balancer, it is possible to implement a network policy mechanism. In GKE, the load balancer does not matter if pods connect using a domain name; it will show the internal IPs of the pods. However, from external sources, it will display the real external IPs (e.g., visitor IPs). By configuring a network policy, you can whitelist only the pod IP range (using CIDR) to allow connections to the MySQL domain name. This ensures that the load balancer dedicated to MySQL remains secure.
> Additionally, keep in mind that the approach may vary by cloud provider. For instance, in AWS or other cloud, you might need to configure everything from scratch for sailing ⛵ ☸.

3. Create a persistent volume claim for MySQL storage:

   ```bash
   kubectl apply -f mysql-storage.yaml
   ```

   Update the `mysql-storage.yaml` file with your desired storage configuration (e.g., size, storage class).

> [!NOTE]
> When storage is `fully encrypted` with the flexibility to be `attached/detached` by the cluster and is bound to the deployment along with VPA,
> the `MySQL data` in the pods will be safe and secure from data loss when pods are restarting (e.g., `restarting 9999999 times`) or other events occur (e.g., `node scaling down from 999999999999 nodes with high specifications (e.g., 999999999999vCPU, 999999999999 Memory Terabyte) to a single node with lower specifications by the autopilot/cluster autoscaler`). This is because all `MySQL data` is bound to the disk.
> Also note that it is not recommended to deploy MySQL along with HPA while using storage that is fully encrypted with the flexibility to be attached/detached.
> This is because the pods will remain in a pending state due to storage limitations, as typically only one pod can access the storage at a time. Even if the storage supports multiple pods (sharing), it is usually limited to a few pods (e.g., `5 pods`).

4. Apply the deployment file to your Kubernetes cluster:

   ```bash
   kubectl apply -f mysql-deploy.yaml
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

> [!NOTE]
> This example shows how `Vertical Pod Autoscaler (VPA)` works properly when the deployment has `PVC/Storage` attached:
```bash
$ kubectl get vpa -n database
NAME        MODE   CPU    MEM         PROVIDED   AGE
mysql-vpa   Auto   100m   874512384   True       3d13h
```

Adjust the VPA configuration in the `mysql-deploy.yaml` file to suit your application's scaling requirements.

## Customization

The provided deployment files are designed to be customizable. You can modify the resource limits, environment variables, and other configurations according to your application's needs.

## Tips

### K8S Network Performance

#### Well-Known Issue When Running on Kubernetes (DigitalOcean)

##### How to Fix the Issue

To resolve the well-known issue when running on Kubernetes with DigitalOcean, modify your service for the Database Load Balancer (after request it) using the following YAML:

```yaml
      service.beta.kubernetes.io/do-loadbalancer-enable-backend-keepalive: "true"
      service.beta.kubernetes.io/do-loadbalancer-enable-proxy-protocol: "true"
      service.beta.kubernetes.io/do-loadbalancer-hostname: db.example.com
      service.beta.kubernetes.io/do-loadbalancer-http-idle-timeout-seconds: "180"
      service.beta.kubernetes.io/do-loadbalancer-size-unit: "1"
      service.beta.kubernetes.io/do-loadbalancer-tls-passthrough: "true"
```

Ensure that the `service.beta.kubernetes.io/do-loadbalancer-hostname` is correctly set for your database, allowing your REST APIs to connect through it.

> [!NOTE]
> If you are using two load balancers—one for the database as a standalone (without NGINX Ingress) and the second for the application—set `service.beta.kubernetes.io/do-loadbalancer-hostname` to `service.beta.kubernetes.io/do-loadbalancer-hostname: db.example.com` for the database. This adjustment will ensure proper connectivity and help prevent "Connection Reset by Peer" errors.

## Cleanup

To remove the deployed resources from your Kubernetes cluster, run the following commands:

```bash
kubectl delete -f mysql-deploy.yaml
```

This will delete the deployment, service, VPA, and persistent volume claim associated with the MySQL service.
