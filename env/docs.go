// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

// Package env provides a constant environment variable for application.
//
// # Setup
//
// Note: For Kubernetes to create secrets, use the bash script "create_k8s_secret.sh" (located in the "k8s-deployment" directory).
// For example, put "DB_DATABASE=databasename", "DB_USERNAME=dbusername", "DB_PASSWORD=yourdbpasswordpogger" in a file named ".env" or whatever any other desired name.
//
// Example ".env" file:
//
//	DB_DATABASE=databasename
//	DB_USERNAME=dbusername
//	DB_PASSWORD=yourdbpasswordpogger
//
// Then run the "create_k8s_secret.sh" script to create the Kubernetes secrets.
// Note: If your Kubernetes cluster already has a built-in Hardware Security Module (HSM), you don't need to use an external secrets mechanism.
//
// # Compatibility
//
// This boilerplate does not support command-line arguments for setup prior to execution (e.g., for configuring path HTTPS/TLS).
// Instead, it primarily relies on environment variables to minimize security risks. While using command-line arguments is possible,
// it can lead to unnecessary complexity due to the large number of variables involved.
package env
