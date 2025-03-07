# EKS Automode NodePools

This repository contains custom resource configurations for defining NodePools in an Amazon EKS cluster with Automode enabled. NodePools allow you to specify the desired characteristics and requirements for the nodes that will be provisioned in your cluster.

## Why NodePools are Used

NodePools provide a way to define the desired properties and requirements for the nodes in your EKS cluster. By creating NodePools, you can specify the instance types, CPU and memory requirements, capacity type (Spot or On-Demand), and other node-related configurations.

Using NodePools allows you to have fine-grained control over the nodes that are provisioned in your cluster. You can create different NodePools for different workload requirements, such as high-performance workloads, general-purpose workloads, or cost-optimized workloads.

## NodePool Configuration

The provided NodePool configuration file (`nodepools/arena-high-performance-pool.yaml`) defines a NodePool for high-performance workloads. Let's go through the key components of the configuration:

### Metadata

- `name`: Specifies the name of the NodePool. Modify this to a meaningful name for your NodePool.

### Spec

- `template.metadata.labels`: Defines labels for the nodes in the NodePool. Modify the `billing-team` label to reflect your team or billing allocation.

- `template.spec.nodeClassRef`: Specifies the NodeClass reference for the NodePool. Modify the `name` field to match the desired NodeClass.

- `template.spec.requirements`: Defines the requirements for the nodes in the NodePool.
  - `karpenter.sh/capacity-type`: Specifies the capacity type for the nodes, allowing both Spot and On-Demand instances.
  - `eks.amazonaws.com/instance-category`: Specifies the instance categories allowed for the nodes (c, m, r, t).
  - `eks.amazonaws.com/instance-family`: Specifies the instance families allowed for the nodes (m5, m5d, c5, c5d, r5, r5d, t3a, m6i, c6i).
  - `eks.amazonaws.com/instance-cpu`: Specifies the allowed CPU configurations for the nodes (1, 2, 4, 8, 16, 32).
  - `eks.amazonaws.com/instance-hypervisor`: Specifies the hypervisor type for the nodes (Nitro).
  - `topology.kubernetes.io/zone`: Specifies the availability zones for the nodes. Modify the values to match your desired zones.
  - `kubernetes.io/arch`: Specifies the architecture for the nodes (amd64).

- `template.spec.terminationGracePeriod`: Specifies the termination grace period for the nodes (30 minutes).

- `disruption`: Defines the disruption budget and consolidation policy for the NodePool.
  - `consolidationPolicy`: Specifies when to consolidate nodes (when empty or underutilized).
  - `consolidateAfter`: Specifies the duration after which to consolidate nodes (1 hour).
  - `budgets`: Defines the disruption budget for the nodes (10% of nodes).

- `limits`: Specifies the resource limits for the NodePool.
  - `cpu`: Specifies the maximum total CPU limit for the nodes (100 CPU units).
  - `memory`: Specifies the maximum total memory limit for the nodes (100Gi).
  - `weight`: Specifies the weight of the NodePool (10).

> [!NOTE]
> Make sure to review and modify the values in the configuration file to match your specific requirements and environment.
