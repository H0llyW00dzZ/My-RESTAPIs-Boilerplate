// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Package archive provides functionality for monitoring and archiving document files based on their size.
//
// # Compatibility
//
// Standard Library Packages:
//
// The package uses the following standard library packages:
//   - [archive/tar]: For creating tar archives.
//   - [compress/gzip]: For compressing the tar archive using gzip.
//   - [fmt]: For string formatting.
//   - [io]: For I/O operations.
//   - [os]: For file and directory operations.
//   - [path/filepath]: For file path manipulation.
//   - [time]: For time-related operations.
//
// File Support:
//
// The package supports archiving any type of file, including but not limited to:
//   - Text files (e.g., .txt, .log)
//   - Binary files (e.g., .bin, .dat)
//   - Compressed files (e.g., .zip, .tar.gz)
//   - Image files (e.g., .jpg, .png)
//   - Document files (e.g., .pdf, .docx)
//
// Storage Support:
//
// The package supports the following storage options:
//   - Local disk storage: The document file and archive directory can be specified as local file paths.
//   - S3 (Bucket) storage: The package can be used with S3 or S3-compatible storage buckets through Kubernetes.
//     The specific storage implementation depends on the CSI driver.
//
// List of supported CSI-Driver S3 or S3-compatible storage (preferred) in Kubernetes:
//
//   - https://github.com/yandex-cloud/k8s-csi-s3 (Tested and works well with own and IaaS on Sandbox)
//   - AWS S3 (Should be supported even without testing on Sandbox)
//
// Note: For S3-compatible storage (preferred), you have the ability to implement your own S3 storage mechanism for infrastructure that doesn't actually use the service, such as IaaS (Infrastructure as Service).
// It's also worth noting that for local disk storage in Kubernetes, it depends on the CSI driver provided by the cloud provider (by attaching external storage not overlay storage such as ephemeral storage).
// However, it's recommended to use a better storage mechanism and driver for reliability, security, and high performance (not slow like HDD).
//
// Supported Multiple Storage in Kubernetes:
//
// In a Kubernetes deployment, you can attach multiple storage options using different CSI drivers or storage classes. This allows for flexibility in data persistence and storage performance, enabling you to choose the appropriate storage solution based on your application requirements.
//
// When attaching multiple storage options to a Kubernetes deployment, consider the following:
//   - Use CSI drivers that integrate with durable storage solutions like cloud-provider-specific block storage (e.g., AWS EBS, GCP Persistent Disk) or distributed file systems (e.g., Ceph, GlusterFS) for reliable and persistent storage.
//   - Choose the appropriate storage class based on the performance characteristics and features required by your application (e.g., SSD vs. HDD, replication, snapshots).
//   - Consider the scalability and cost implications of using multiple storage options, as different storage solutions may have different pricing models and scalability limitations.
//
// How Multiple Storage in Kubernetes Works:
//
// When you attach external storage for multiple storage options, it is mounted from the path where it is bound (e.g., /home/storage1, /home/storage2). When multiple storage options are attached, avoid setting the path with "." (e.g., ./home/storage1).
// If you use "." when you already have multiple storage options attached, it will use the storage overlay (e.g., ephemeral storage) from the container instead of the actual multiple storage options that are attached.
//
// Integration with Fiber Middleware Logs:
//
// The package can be seamlessly integrated with Fiber middleware logs for archiving purposes. It supports archiving Fiber middleware logs stored in the following locations:
//   - Local disk: Fiber middleware logs stored on the local disk can be directly archived using this package.
//   - S3 (Bucket) or S3-compatible (Bucket) storage through Kubernetes: If Fiber middleware logs are stored in an S3 bucket through Kubernetes (depending on the CSI driver), the package can be configured to archive those logs.
//   - S3 (Bucket) or S3-compatible (Bucket) storage without Kubernetes: The package can also be used to archive Fiber middleware logs stored in an S3 bucket without Kubernetes, by specifying the appropriate S3 configuration.
//
// Deployment Type:
//
// The package is primarily designed to be stable with Vertical Pod Autoscaling (VPA) rather than Horizontal Pod Autoscaling (HPA).
// It is recommended to use VPA for scaling the deployment based on resource requirements.
//
// Note that Regarding the deployment type, it is generally considered a bad practice to set it to "Stateful" even when external storage or multiple external storage options are attached.
// This applies to various roles such as Developers, DevOps, DevSecOps, or any other similar positions, unless you are specifically focusing on Kubernetes components like drivers or other specialized areas.
//
// Even for game servers running on Kubernetes using this boilerplate for interacting with the cluster and setting up the game server, it is not recommended to set the deployment type to "Stateful".
// In Kubernetes, it is possible to run game servers and other deployments, such as AI, without using stateful deployments.
//
// The reason behind this recommendation is that stateful deployments introduce additional complexity and management overhead compared to stateless deployments.
// Stateless deployments are generally more scalable, easier to manage, and provide better flexibility in terms of resource allocation and scaling.
//
// However, there may be specific scenarios where stateful deployments are necessary, such as when dealing with persistent data that requires strict consistency and ordering guarantees.
// It's worth noting that even for databases like MySQL, it is still possible to run them stably as stateless deployments by attaching external storage + VPA. Such cases are relatively rare.
// If stateful deployments are required, careful consideration and design are necessary to ensure the proper handling of stateful components within the Kubernetes environment.
//
// It's important to evaluate the specific requirements and characteristics of your application and determine the most appropriate deployment strategy based on those factors.
//
// # Security Considerations
//
// In Kubernetes, the security risk is relatively low because you have control over the permissions. However, if your deployment uses a minimal image that does not fully interact with the operating system,
// it is generally safe to run as root, as long as the image is minimal and has limited capabilities.
//
// When running containers with minimal images, consider the following security best practices:
//
//   - Use a minimal base image that includes only the necessary dependencies and libraries required by your application.
//   - Avoid installing unnecessary packages or tools that could potentially introduce security vulnerabilities.
//   - Ensure that the container runtime has limited access to the host system resources and follows the principle of least privilege.
//   - Regularly scan and update the minimal image for any known security vulnerabilities and apply the latest security patches.
//   - Implement proper network segmentation and access controls to limit the potential impact of a compromised container.
//
// It's important to note that running containers as root should still be approached with caution and only when absolutely necessary.
// Whenever possible, it is recommended to run containers with a non-root user and grant only the specific permissions required by the application.
//
// By following these security considerations and best practices, you can mitigate the risks associated with running containers as root in a minimal image environment.
//
// # Usage
//
// The main entry point of the package is the Do function, which continuously monitors a specified
// document file and archives it when its size reaches the configured maximum size. The archiving
// process involves compressing the document file into a tar.gz archive and storing it in the
// specified archive directory.
//
// The package provides a Config struct to configure the archiving process, including the maximum
// file size and the check interval for monitoring the file size. The DefaultConfig function returns
// a Config instance with default values.
//
// Note: The package assumes that the caller has the necessary permissions to read the document file
// and write to the archive directory.
package archive
