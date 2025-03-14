# EKS Automode NodePools

This repository contains custom resource configurations for defining NodePools in an Amazon EKS cluster with Automode enabled. NodePools allow you to specify the desired characteristics and requirements for the nodes that will be provisioned in your cluster.

## Why NodePools are Used

NodePools provide a way to define the desired properties and requirements for the nodes in your EKS cluster. By creating NodePools, you can specify the instance types, CPU and memory requirements, capacity type (Spot or On-Demand), and other node-related configurations.

Using NodePools allows you to have fine-grained control over the nodes that are provisioned in your cluster. You can create different NodePools for different workload requirements, such as high-performance workloads, general-purpose workloads, or cost-optimized workloads.

## NodePool Configurations

This repository contains two NodePool configuration files:

1. `nodepools/arena-high-performance-pool.yaml`: Defines a NodePool for high-performance workloads.
2. `nodepools/critical-pool.yaml`: Defines a NodePool for critical workloads.

Let's go through the key components of each configuration:

### Arena High-Performance Pool

The `arena-high-performance-pool.yaml` file defines a NodePool for high-performance workloads. It includes the following key components:

- `metadata.name`: Specifies the name of the NodePool. Modify this to a meaningful name for your high-performance NodePool.
- `template.metadata.labels`: Defines labels for the nodes in the NodePool. Modify the `billing-team` label to reflect your team or billing allocation.
- `template.spec.nodeClassRef`: Specifies the NodeClass reference for the NodePool. Modify the `name` field to match the desired NodeClass.
- `template.spec.requirements`: Defines the requirements for the nodes in the NodePool, including capacity type, instance categories, instance families, CPU configurations, hypervisor type, availability zones, and architecture.
- `template.spec.terminationGracePeriod`: Specifies the termination grace period for the nodes (1 hour).
- `disruption`: Defines the disruption budget and consolidation policy for the NodePool.
- `limits`: Specifies the resource limits for the NodePool, including CPU, memory, and weight.

### Critical Pool

The `critical-pool.yaml` file defines a NodePool for critical workloads. It includes the following key components:

- `metadata.name`: Specifies the name of the NodePool. Ensure that the name reflects the purpose of this critical NodePool.
- `template.metadata.labels`: Defines labels for the nodes in the NodePool. Adjust the `billing-team` label to match your organization's conventions.
- `template.spec.nodeClassRef`: Specifies the NodeClass reference for the NodePool. Make sure the `name` field matches the corresponding NodeClass resource.
- `template.spec.requirements`: Defines the requirements for the nodes in the NodePool, including capacity type, instance categories, instance families, CPU configurations, hypervisor type, availability zones, and architecture.
- `template.spec.taints`: Applies the "CriticalAddonsOnly" taint to ensure that only critical pods are scheduled on these nodes.
- `template.spec.terminationGracePeriod`: Specifies the termination grace period for the nodes (1 hour 30 minutes).
- `disruption`: Defines the consolidation policy and interval for the NodePool.
- `limits`: Specifies the resource limits for the NodePool, including CPU, memory, and weight.

> [!NOTE]
> Make sure to review and modify the values in the configuration files to match your specific requirements and environment.

## Best Practices for Controlling Node Pools via EKS Automode

1. Ensure that the node pools match your specific requirements and environment. For example, the `critical-pool` is suitable for deploying Kubernetes components such as CoreDNS with HPA, Ingress Nginx, and other deployments that rely on the `kube-system` namespace.

2. Use meaningful names for your node pools that reflect their purpose and workload characteristics. This helps in identifying and managing node pools effectively.

3. Adjust the labels, node class references, and resource requirements according to your organization's conventions and the specific needs of your workloads.

4. Consider using Spot instances for cost optimization, especially for non-critical workloads. However, ensure that you have appropriate fallback mechanisms in place.

5. Utilize taints and tolerations to control pod scheduling on specific node pools. For example, the `critical-pool` uses the "CriticalAddonsOnly" taint to ensure that only critical pods are scheduled on those nodes.

6. Set appropriate termination grace periods for your node pools based on the nature of your workloads and the time required for graceful shutdown.

7. Define disruption budgets and consolidation policies to minimize the impact of node terminations and optimize resource utilization.

8. Regularly review and adjust the resource limits (CPU, memory, weight) of your node pools based on actual usage patterns and the evolving needs of your workloads.

By following these best practices and carefully configuring your node pools, you can ensure optimal performance, reliability, and cost-efficiency for your EKS cluster.

> [!NOTE]
> The current custom resource configurations for the EKS Automode NodePools include both `spot` and `on-demand` capacity types. It's important to note that using `spot` instances may not guarantee zero downtime due to the nature of spot instances, which can be interrupted or terminated by AWS based on market conditions. To achieve zero downtime, it is recommended to remove the `spot` capacity type from the `arena-high-performance-pool.yaml` and `critical-pool.yaml` configuration files and rely solely on `on-demand` instances.
