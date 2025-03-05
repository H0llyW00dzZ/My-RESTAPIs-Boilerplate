# AWS EBS StorageClass Configurations

This repository contains StorageClass configurations for provisioning AWS EBS volumes in a Kubernetes cluster. When creating a new cluster in AWS, you need to set up various resources manually, including StorageClasses. These StorageClass definitions provide a convenient way to provision EBS volumes with different configurations.

## Why StorageClasses are Added

When you create a new Kubernetes cluster in AWS, it doesn't come with pre-configured StorageClasses for EBS volumes. By default, you would have to manually create and configure StorageClasses to provision EBS volumes with desired settings such as volume type, encryption, and reclaim policy.

Adding these StorageClass definitions to your cluster ensures that you have a set of pre-defined configurations ready to use. This saves time and effort in manually creating StorageClasses each time you set up a new cluster.

## StorageClass Configurations

The provided StorageClass configurations cover different scenarios and requirements:

1. `gp2-encrypted-retain`: Provisions encrypted gp2 volumes with the `Retain` reclaim policy.
2. `gp2-retain`: Provisions non-encrypted gp2 volumes with the `Retain` reclaim policy.
3. `gp2-encrypted`: Provisions encrypted gp2 volumes with the `Delete` reclaim policy.
4. `gp2`: Provisions non-encrypted gp2 volumes with the `Delete` reclaim policy.
5. `gp3-encrypted-retain`: Provisions encrypted gp3 volumes with the `Retain` reclaim policy.
6. `gp3-retain`: Provisions non-encrypted gp3 volumes with the `Retain` reclaim policy.
7. `gp3-encrypted`: Provisions encrypted gp3 volumes with the `Delete` reclaim policy.
8. `gp3`: Provisions non-encrypted gp3 volumes with the `Delete` reclaim policy.

These StorageClasses cover both gp2 and gp3 volume types, with options for encryption and different reclaim policies (`Retain` and `Delete`).

> [!NOTE]
> Make sure the CSI Driver for AWS EBS is installed via add-ons before applying these StorageClasses.
