# AWS EBS StorageClass Configurations

This repository contains StorageClass configurations for provisioning AWS EBS volumes in a Kubernetes cluster. When creating a new cluster in AWS, you need to set up various resources manually, including StorageClasses. These StorageClass definitions provide a convenient way to provision EBS volumes with different configurations.

## Why StorageClasses are Added

When you create a new Kubernetes cluster in AWS, it doesn't come with pre-configured StorageClasses for EBS volumes. By default, you would have to manually create and configure StorageClasses to provision EBS volumes with desired settings such as volume type, encryption, and reclaim policy.

Adding these StorageClass definitions to your cluster ensures that you have a set of pre-defined configurations ready to use. This saves time and effort in manually creating StorageClasses each time you set up a new cluster.

## StorageClass Configurations

The provided StorageClass configurations cover different scenarios and requirements:

### General Purpose SSD (gp2 and gp3) StorageClasses

File location: `ebs/gp.yaml`

1. `gp2-encrypted-retain`: Provisions encrypted gp2 volumes with the `Retain` reclaim policy.
2. `gp2-retain`: Provisions non-encrypted gp2 volumes with the `Retain` reclaim policy.
3. `gp2-encrypted`: Provisions encrypted gp2 volumes with the `Delete` reclaim policy.
4. `gp2`: Provisions non-encrypted gp2 volumes with the `Delete` reclaim policy.
5. `gp3-encrypted-retain`: Provisions encrypted gp3 volumes with the `Retain` reclaim policy.
6. `gp3-retain`: Provisions non-encrypted gp3 volumes with the `Retain` reclaim policy.
7. `gp3-encrypted`: Provisions encrypted gp3 volumes with the `Delete` reclaim policy.
8. `gp3`: Provisions non-encrypted gp3 volumes with the `Delete` reclaim policy.

These StorageClasses cover both gp2 and gp3 volume types, with options for encryption and different reclaim policies (`Retain` and `Delete`).

### Provisioned IOPS SSD (io1 and io2) StorageClasses

File location: `ebs/io.yaml`

1. `io1-encrypted-retain`: Provisions encrypted io1 volumes with the `Retain` reclaim policy and 50 IOPS per GiB.
2. `io1-retain`: Provisions non-encrypted io1 volumes with the `Retain` reclaim policy and 50 IOPS per GiB.
3. `io1-encrypted`: Provisions encrypted io1 volumes with the `Delete` reclaim policy and 50 IOPS per GiB.
4. `io1`: Provisions non-encrypted io1 volumes with the `Delete` reclaim policy and 50 IOPS per GiB.
5. `io2-encrypted-retain`: Provisions encrypted io2 volumes with the `Retain` reclaim policy and 500 IOPS per GiB.
6. `io2-retain`: Provisions non-encrypted io2 volumes with the `Retain` reclaim policy and 500 IOPS per GiB.
7. `io2-encrypted`: Provisions encrypted io2 volumes with the `Delete` reclaim policy and 500 IOPS per GiB.
8. `io2`: Provisions non-encrypted io2 volumes with the `Delete` reclaim policy and 500 IOPS per GiB.

These StorageClasses cover both io1 and io2 volume types, with options for encryption, different reclaim policies (`Retain` and `Delete`), and specified IOPS per GiB.

> [!NOTE]
> Make sure the CSI Driver for AWS EBS is installed via add-ons before applying these StorageClasses.

## StorageClass Configurations for EKS Automode

In addition to the standard StorageClass configurations, this repository also includes StorageClass configurations specifically tailored for Amazon EKS clusters with automode enabled.

### General Purpose SSD (gp2 and gp3) StorageClasses for EKS Automode

File location: `auto/ebs/gp.yaml`

1. `gp2-auto-encrypted`: Automatically provisions encrypted gp2 volumes.
2. `gp2-auto`: Automatically provisions non-encrypted gp2 volumes.
3. `gp3-auto-encrypted`: Automatically provisions encrypted gp3 volumes.
4. `gp3-auto`: Automatically provisions non-encrypted gp3 volumes.

These StorageClasses are designed to work seamlessly with EKS automode, allowing automatic provisioning of EBS volumes with gp2 and gp3 volume types, with options for encryption.

### Provisioned IOPS SSD (io1 and io2) StorageClasses for EKS Automode

File location: `auto/ebs/io.yaml`

1. `io1-auto-encrypted`: Automatically provisions encrypted io1 volumes with 50 IOPS per GiB.
2. `io1-auto`: Automatically provisions non-encrypted io1 volumes with 50 IOPS per GiB.
3. `io2-auto-encrypted`: Automatically provisions encrypted io2 volumes with 500 IOPS per GiB.
4. `io2-auto`: Automatically provisions non-encrypted io2 volumes with 500 IOPS per GiB.

These StorageClasses are designed to work seamlessly with EKS automode, allowing automatic provisioning of EBS volumes with io1 and io2 volume types, with options for encryption and specified IOPS per GiB.
