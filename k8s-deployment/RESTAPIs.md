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
> Since this deployment supports `100% HPA`, which is suitable for handling billions of requests/workers (combined with [worker package](https://github.com/H0llyW00dzZ/My-RESTAPIs-Boilerplate/tree/master/worker)) `concurrently and efficiently`, it is recommended not to attach it with any `storage (PVC, PV)` to this deployment.
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

> [!NOTE]
> This example shows how `Horizontal Pod Autoscaler (HPA)` works properly, handling billions of requests/workers (combined with the [worker package](https://github.com/H0llyW00dzZ/My-RESTAPIs-Boilerplate/tree/master/worker)) `concurrently and efficiently`:

- **Events:**

```
12m         Normal   SuccessfulRescale   horizontalpodautoscaler/senior-golang-worker-hpa   New size: 14; reason: cpu resource utilization (percentage of request) below target
58m         Normal   SuccessfulRescale   horizontalpodautoscaler/senior-golang-worker-hpa   New size: 17; reason: cpu resource utilization (percentage of request) above target
12m         Normal   ScalingReplicaSet   deployment/senior-golang-worker                    Scaled down replica set senior-golang-worker-84bcb968 to 14 from 17
58m         Normal   ScalingReplicaSet   deployment/senior-golang-worker                    Scaled up replica set senior-golang-worker-84bcb968 to 17 from 14
```

- **Describe HPA:**

```
Name:                                                  senior-golang-worker-hpa
Namespace:                                             senior-golang-worker
Labels:                                                <none>
Annotations:                                           <none>
CreationTimestamp:                                     Wed, 11 Sep 2024 19:44:05 +0700
Reference:                                             Deployment/senior-golang-worker
Metrics:                                               ( current / target )
  resource cpu on pods  (as a percentage of request):  76% (268m) / 80%
Min replicas:                                          1
Max replicas:                                          30
Deployment pods:                                       14 current / 14 desired
Conditions:
  Type            Status  Reason              Message
  ----            ------  ------              -------
  AbleToScale     True    ReadyForNewScale    recommended size matches current size
  ScalingActive   True    ValidMetricFound    the HPA was able to successfully calculate a replica count from cpu resource utilization (percentage of request)
  ScalingLimited  False   DesiredWithinRange  the desired count is within the acceptable range
Events:
  Type    Reason             Age                   From                       Message
  ----    ------             ----                  ----                       -------
  Normal  SuccessfulRescale  60m (x5 over 34h)     horizontal-pod-autoscaler  New size: 17; reason: cpu resource utilization (percentage of request) above target
  Normal  SuccessfulRescale  14m (x12 over 2d23h)  horizontal-pod-autoscaler  New size: 14; reason: cpu resource utilization (percentage of request) below target
```

- **Watching HPA (Stable for long-running (Smooth Sailing ⛵ ☸) processes in combination with [`worker package`](https://github.com/H0llyW00dzZ/My-RESTAPIs-Boilerplate/tree/master/worker)):**

```
b0zal@Linux:~$ kubectl get hpa --watch
        NAME                           REFERENCE                TARGETS    MINPODS   MAXPODS   REPLICAS    AGE
senior-golang-worker-hpa   Deployment/senior-golang-worker   cpu: 75%/80%   1         30        22         18d
senior-golang-worker-hpa   Deployment/senior-golang-worker   cpu: 78%/80%   1         30        22         18d
senior-golang-worker-hpa   Deployment/senior-golang-worker   cpu: 68%/80%   1         30        22         18d
senior-golang-worker-hpa   Deployment/senior-golang-worker   cpu: 66%/80%   1         30        22         18d
senior-golang-worker-hpa   Deployment/senior-golang-worker   cpu: 54%/80%   1         30        22         18d
senior-golang-worker-hpa   Deployment/senior-golang-worker   cpu: 64%/80%   1         30        22         18d
senior-golang-worker-hpa   Deployment/senior-golang-worker   cpu: 34%/80%   1         30        22         18d
senior-golang-worker-hpa   Deployment/senior-golang-worker   cpu: 71%/80%   1         30        22         18d
senior-golang-worker-hpa   Deployment/senior-golang-worker   cpu: 76%/80%   1         30        22         18d
senior-golang-worker-hpa   Deployment/senior-golang-worker   cpu: 75%/80%   1         30        22         18d
senior-golang-worker-hpa   Deployment/senior-golang-worker   cpu: 65%/80%   1         30        22         18d
senior-golang-worker-hpa   Deployment/senior-golang-worker   cpu: 68%/80%   1         30        22         18d
senior-golang-worker-hpa   Deployment/senior-golang-worker   cpu: 61%/80%   1         30        22         18d
senior-golang-worker-hpa   Deployment/senior-golang-worker   cpu: 57%/80%   1         30        22         18d
senior-golang-worker-hpa   Deployment/senior-golang-worker   cpu: 25%/80%   1         30        22         18d
senior-golang-worker-hpa   Deployment/senior-golang-worker   cpu: 17%/80%   1         30        22         18d
senior-golang-worker-hpa   Deployment/senior-golang-worker   cpu: 53%/80%   1         30        22         18d
senior-golang-worker-hpa   Deployment/senior-golang-worker   cpu: 76%/80%   1         30        22         18d
```

The `cpu: 78%/80%` going `up/down` and `REPLICAS 22` are `dynamic` (depending on how many tasks/jobs are processed by the [`worker package`](https://github.com/H0llyW00dzZ/My-RESTAPIs-Boilerplate/tree/master/worker)) and not caused by the `cluster` itself or the `kernel`. It is suitable for [cluster autoscaler](https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler) (auto-pilot).

> [!WARNING]  
> There is also a `warning` regarding the [cluster autoscaler](https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler) (auto-pilot). When using the [CA](https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler) (auto-pilot) in the cluster, consider not explicitly specifying [`node selectors`](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodeselector) for each `deployment/pod`, as it can make it difficult for the [CA](https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler) (auto-pilot) to maintain the nodes (e.g., `self-healing`, `detach nodes when no longer needed`).

> [!NOTE]
> It's important to note that `Horizontal Pod Autoscaler (HPA)` can be used with various types of software and applications, including `websites` and `game servers`, as long as they are deployed as pods in a Kubernetes cluster. The choice between `HPA` and `Vertical Pod Autoscaler (VPA)` depends on the specific requirements and characteristics of the workload. For example, for a `game server` [Counter-Strike 2](https://www.counter-strike.net/cs2) `Community Servers`, the stability and performance may depend on factors such as server hardware, network infrastructure, and configuration (if you have a deep understanding of Kubernetes, this can be easily managed), rather than solely on the use of `HPA` or `VPA`. Based on personal experience hosting a [Counter-Strike 2](https://www.counter-strike.net/cs2) `game server` fully managed and isolated by `Kubernetes`, it was found to be more stable than the official servers provided by `Steam` or `Faceit`.

Adjust the HPA configuration in the `restapis-deploy.yaml` file to suit your application's scaling requirements.

- **Average CPU and Memory Usage When Idle (as viewed on Grafana):**

<p align="center">
   <img src="https://i.imgur.com/y9Ky3xZ.png" alt="fully-managed-and-isolated-by-k8s">
</p>

- **Average CPU (on AMD) and Memory Usage During Smooth Sailing ⛵ ☸ in Combination with the [`worker package`](https://github.com/H0llyW00dzZ/My-RESTAPIs-Boilerplate/tree/master/worker) (as viewed on Grafana):**

<p align="center">
   <img src="https://i.imgur.com/4u7i9VV.png" alt="fully-managed-and-isolated-by-k8s">
   <img src="https://i.imgur.com/x5oRx7h.png" alt="fully-managed-and-isolated-by-k8s">
</p>

As you can see, the memory usage is dynamic yet `stable and predictable`, unlike static memory usage where CPU growth affects all memory uniformly. This ensures smooth sailing ⛵ ☸, thanks to a new Go package called [`Unique`](https://pkg.go.dev/unique).

## Customization

The provided deployment files are designed to be customizable. You can modify the resource limits, environment variables, and other configurations according to your application's needs. Additionally, you can adjust the Ingress configuration to match your desired routing rules and TLS settings.

## Tips

### K8S Network Performance

Here are some tips to boost/improve the network performance. These are well-known in GKE because most important components are already built-in, and you only need to enable them:

- GKE:

Enable [Dataplane V2 Observability](https://docs.cilium.io/en/stable/internals/hubble/)

> [!NOTE]
> It can also improve the security mechanism if you have a deep understanding of networking.
> Other cloud providers might be added later, as I don't have experience with them. Additionally, the K8s ecosystem is large, not small.

### REST API Concurrency with HPA

- **Install [Cluster Autoscaler (CA)](https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler)**: If you are using **GKE Autopilot**, you do not need to install the CA manually, as it is managed for you.

To enhance REST API concurrency (in combination with the [`worker package`](https://github.com/H0llyW00dzZ/My-RESTAPIs-Boilerplate/tree/master/worker)) and improve HPA performance:

- Start with a small deployment (e.g., **1 node**), and set the **maximum CPU** to **350m** with **80% utilization** for the HPA. This strategy allows for scaling as demand increases.

#### Example of REST API Concurrency with HPA

```bash
b0zal@Linux:~$ kubectl get hpa
NAME                        REFERENCE                   TARGETS        MINPODS   MAXPODS   REPLICAS    AGE
senior-golang-worker-hpa   Deployment/senior-golang   cpu: 75%/80%      1         60        41         35d
```

> [!NOTE]
> In this example, the REST API ([This Repo](https://github.com/H0llyW00dzZ/My-RESTAPIs-Boilerplate)) is running with **22 vCPUs** (extremely scalable) across **8 nodes**, for **41 Pods**.

### Example of 80% Utilization for the HPA in Math ([`LaTeX`](https://en.wikipedia.org/wiki/LaTeX))

To understand how the utilization and vCPUs work together, consider the following calculations:

1. **Calculate total vCPUs used by Pods:**
   - Each pod uses **350m** (which is **0.35 vCPUs**).
   - With **41 Pods**, the total vCPUs used is:

$$
   \text{Total vCPUs} = 41 \text{ Pods} \times 0.35 \text{ vCPUs/Pod} = 14.35 \text{ vCPUs}
$$

2. **Calculate required vCPUs based on utilization:**
   - If you are targeting **80% utilization**, the required vCPUs can be calculated as:

$$
   \text{Required vCPUs} = \frac{\text{Total vCPUs}}{\text{Utilization}} = \frac{14.35}{0.80} \approx 17.94 \text{ vCPUs}
$$

This means that with **22 vCPUs** available, you have sufficient capacity to handle the load of **41 Pods** while maintaining an 80% utilization target.

### Well-Known Issue: `Connection Reset by Peer` When Running on Kubernetes (DigitalOcean)

##### How to Fix the Issue

To resolve the well-known issue `Connection Reset by Peer` when running on Kubernetes with DigitalOcean, modify your service for the NGINX Ingress (after installing it) using the following YAML:

```yaml
      meta.helm.sh/release-name: ingress-nginx
      meta.helm.sh/release-namespace: ingress-nginx
      service.beta.kubernetes.io/do-loadbalancer-enable-backend-keepalive: "true"
      service.beta.kubernetes.io/do-loadbalancer-enable-proxy-protocol: "true"
      service.beta.kubernetes.io/do-loadbalancer-hostname: api.example.com
      service.beta.kubernetes.io/do-loadbalancer-http-idle-timeout-seconds: "180"
      service.beta.kubernetes.io/do-loadbalancer-size-unit: "1"
      service.beta.kubernetes.io/do-loadbalancer-tls-passthrough: "true"
```

Make sure to modify the `service.beta.kubernetes.io/do-loadbalancer-hostname` within your REST APIs.

> [!NOTE]
> If you are using two load balancers (one for the database as a standalone without NGINX Ingress, and the second for the application), change `service.beta.kubernetes.io/do-loadbalancer-hostname` to `service.beta.kubernetes.io/do-loadbalancer-hostname: db.example.com` for the database. This will ensure proper connectivity and prevent "Connection Reset by Peer" errors.
>
> Additionally, the `Connection Reset by Peer` error can occur when pods cannot communicate with each other or with themselves. For example, if your pod's IP is `10.0.0.1` and you try to use `curl` to access it via `example.com`, which is bound to `10.0.0.1`, you may encounter the `Connection Reset by Peer` error. However, using `curl` directly to `10.0.0.1` would work properly. This issue can arise even within the same `virtual machine`.

### Setup DOKS External Load Balancer Hostname for Ingress-NGINX

To set up a DOKS external load balancer hostname that allows pods to communicate and prevents the error `Connection Reset by Peer`, follow these steps. This setup enables any domain based on DNS.

1. Ensure you replace `service.beta.kubernetes.io/do-loadbalancer-hostname: api.example.com` with `service.beta.kubernetes.io/do-loadbalancer-hostname: host.example.com`:
   ```yaml
   service.beta.kubernetes.io/do-loadbalancer-hostname: host.example.com
   ```

2. In `host.example.com`, set the [DNS A record](https://www.cloudflare.com/learning/dns/dns-records/dns-a-record/) with the IP of the ingress load balancer to `host.example.com`.

3. When creating multiple ingresses across different services, you don't need to set each domain to an IP address. Instead, use a [CNAME](https://www.cloudflare.com/learning/dns/dns-records/dns-cname-record/) record pointing to `host.example.com`.

> [!TIP]
> To enhance the hostname, you can create a random name, for example, using a SHA-256 digest, such as `936a185caaa266bb9cbe981e9e05cb78cd732b0b3280eb944412bb6f8f8f07af.example.com`.

## Cleanup

To remove the deployed resources from your Kubernetes cluster, run the following commands:

```bash
kubectl delete -f restapis-ingress.yaml
kubectl delete -f restapis-deploy.yaml
```

This will delete the Ingress resource, deployment, service, and HPA associated with the REST API application.
