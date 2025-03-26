# Gitea Docker-in-Docker (DinD) Runner

<p align="center">
   <img src="https://kubernetes.io/images/kubernetes.png" alt="sailing-with-k8s" width="80">
   <img src="https://i.imgur.com/wGetVaj.png" alt="The-Black-Pearl" width="80">
</p>

This repository contains the configuration for deploying a Gitea runner using Docker-in-Docker (DinD) in a Kubernetes environment.

## Overview

The Gitea DinD Runner allows you to execute CI/CD jobs within a Kubernetes cluster using Docker containers. It leverages the act-runner to manage job execution and supports caching and custom network configurations.

> [!NOTE]
> The current deployment for the Gitea DinD Runner is compatible with the existing workflow implementation in [`private-registry-kvm-only.yaml`](https://github.com/H0llyW00dzZ/My-RESTAPIs-Boilerplate/blob/master/.github/workflows/private-registry-kvm-only.yaml).
>
> For Example:
> ![Screenshot from 2025-03-18 23-53-13](https://i.imgur.com/yN0wCBQ.png)
>
> ![Screenshot from 2025-03-20 05-12-57](https://i.imgur.com/gcIpDFy.png)
> 
> ![Screenshot from 2025-03-18 23-53-57](https://i.imgur.com/vtteJ6g.png)
>
> Note that Gitea itself (self-hosted) is also running in a Kubernetes container.
>
> The performance is also different when building images with QEMU (e.g., it won't get stuck forever when building multi-arch images, unlike building outside of the cluster, which can cause the process to get stuck forever or become very slow).
> This is especially true when your cluster has built-in custom resources such as node pools like EKS Auto Mode. It basically becomes a layer on top of Kubernetes, because you can specify node specs through node pools, and your infrastructure always runs smoothly while sailing 🛳️ the Kubernetes seas ☸.
>
> The Docker context also uses TCP with TLS enabled, rather than the default socket (e.g., `unix:///var/run/docker.sock`), to ensure compatibility with DinD (Docker-in-Docker) managed by Kubernetes.

## Prerequisites

- Kubernetes cluster (v1.32.0 or later)
- Gitea instance
- Docker registry access
- Persistent storage (optional, e.g., EFS for RWX)

## Deployment

### Configuration

The configuration is managed using a `ConfigMap` and a `Deployment` in Kubernetes.

1. **ConfigMap**: Defines the runner settings, including logging, environment variables, and caching options.
2. **Deployment**: Manages the runner pods, configuring resources and mounting necessary volumes.

### Steps

1. **Create a Namespace**

   ```bash
   kubectl create namespace gitea
   ```

2. **Apply ConfigMap & Deployment**

   Update the deployment file (`gitea-runner-dind.yaml`) with your specific settings, particularly the `GITEA_INSTANCE_URL` and `GITEA_RUNNER_REGISTRATION_TOKEN`. Once you have verified that the configuration is correct, apply the changes using the following command:

   ```bash
   kubectl apply -f gitea-runner-dind.yaml
   ```

## Configuration Details

### Logging

- **Level**: Set the desired logging level (`trace`, `debug`, `info`, `warn`, `error`, `fatal`).

### Runner Settings

- **Capacity**: Number of concurrent tasks.
- **Environment Variables**: Define custom environment variables for jobs.
- **Timeouts**: Configure job and shutdown timeouts.
- **Labels**: Specify labels to determine job execution criteria.

### Caching

- **Enable**: Toggle cache server usage.
- **Directory**: Specify where cache data is stored.
- **External Server**: Use an external cache server if needed.

### Container Options

- **Network**: Define the network configuration for containers.
- **Privileged Mode**: Required for Docker-in-Docker operations.
- **Volumes**: Configure valid volumes for mounting.

## Tips

- For storage solutions like EFS or RWX (Storage Classes that support RWX), you can increase the number of replicas for scalability.
- Ensure all placeholder values are replaced with actual configurations before deployment.
- If your Kubernetes cluster uses custom resources, such as EKS Auto Mode, it will become easier & effectively manage node pools for the runner in your infrastructure.

## Security Considerations

- **TLS Verification**: Ensure `insecure` is set to `false` for production environments.
- **Secrets Management**: Use Kubernetes secrets to manage sensitive information like tokens.
