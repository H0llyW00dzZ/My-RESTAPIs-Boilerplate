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
