# Cluster Autoscaler on AWS via Helm

This README provides instructions on how to deploy the Cluster Autoscaler on an AWS Kubernetes cluster using Helm.

## Prerequisites

Before deploying the Cluster Autoscaler, ensure that you have the following prerequisites:

- An AWS Kubernetes cluster (e.g., Amazon EKS)
- Helm installed on your local machine
- AWS CLI configured with the necessary permissions

## Configuration

1. Create an IAM policy for the Cluster Autoscaler:

   Create an IAM policy using the provided `ca-policy.json` file. Attach this policy to the IAM role associated with your Kubernetes cluster nodes.

2. Configure the Helm values:

   Update the `helm-values.yaml` file with your specific cluster configuration, such as the cluster name, AWS region, and IAM role ARN.

## Deployment

1. Add the Cluster Autoscaler Helm repository:

   ```bash
   helm repo add autoscaler https://kubernetes.github.io/autoscaler
   ```

2. Update the Helm repository:

   ```bash
   helm repo update
   ```

3. Install the Cluster Autoscaler using Helm:

   ```bash
   helm install cluster-autoscaler autoscaler/cluster-autoscaler --namespace kube-system --values helm-values.yaml
   ```

   This command installs the Cluster Autoscaler in the `kube-system` namespace using the provided `helm-values.yaml` file.

4. Verify the installation:

   ```bash
   kubectl get pods -n kube-system | grep cluster-autoscaler
   ```

   You should see the Cluster Autoscaler pod running in the `kube-system` namespace.

> [!NOTE]
> If the installation fails after applying the `helm-values.yaml` file due to RBAC-related issues (e.g., forbidden access), try applying the RBAC configuration manually using the provided `rbac-autoscaler.yaml` file:

```bash
kubectl apply -f rbac-autoscaler.yaml
```

Then, retry the Helm installation command.

## Configuration

The Cluster Autoscaler configuration can be customized by modifying the `helm-values.yaml` file. Some important configuration options include:

- `autoDiscovery.clusterName`: The name of your Kubernetes cluster.
- `awsRegion`: The AWS region where your cluster is deployed.
- `rbac.serviceAccount.annotations.eks.amazonaws.com/role-arn`: The ARN of the IAM role associated with your cluster nodes.
- `extraArgs`: Additional arguments passed to the Cluster Autoscaler, such as `skip-nodes-with-local-storage`, `expander`, `scale-down-unneeded-time`, and `scale-down-delay-after-add`.

For more information on the available configuration options, refer to the [Cluster Autoscaler Helm chart documentation](https://github.com/kubernetes/autoscaler/tree/master/charts/cluster-autoscaler).

## Uninstallation

To uninstall the Cluster Autoscaler, run the following command:

```bash
helm uninstall cluster-autoscaler --namespace kube-system
```

This command removes the Cluster Autoscaler from your Kubernetes cluster.

## Troubleshooting

If you encounter any issues during the deployment or operation of the Cluster Autoscaler, you can check the logs of the Cluster Autoscaler pod using the following command:

```bash
kubectl logs -f deployment/cluster-autoscaler -n kube-system
```

This command displays the logs of the Cluster Autoscaler pod, which can help in identifying and resolving any issues.

For more information and advanced configuration options, refer to the [Cluster Autoscaler documentation](https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler).
