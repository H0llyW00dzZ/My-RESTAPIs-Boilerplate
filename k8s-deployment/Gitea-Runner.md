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

   Ensure the configuration in `gitea-runner-dind.yaml` is correct, then apply:

   ```bash
   kubectl apply -f gitea-runner-dind.yaml
   ```

3. **Deploy the Runner**

   Update the deployment file with your specific settings, particularly the `GITEA_INSTANCE_URL` and `GITEA_RUNNER_REGISTRATION_TOKEN`, then apply:

   ```bash
   kubectl apply -f deployment.yml
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

- For storage solutions like EFS, you can increase the number of replicas for scalability.
- Ensure all placeholder values are replaced with actual configurations before deployment.

## Security Considerations

- **TLS Verification**: Ensure `insecure` is set to `false` for production environments.
- **Secrets Management**: Use Kubernetes secrets to manage sensitive information like tokens.
