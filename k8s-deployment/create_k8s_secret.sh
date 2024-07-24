#!/bin/bash
# Copyright (c) 2024 H0llyW00dz All rights reserved.
# 
# By accessing or using this software, you agree to be bound by the terms
# of the License Agreement, which you can find at LICENSE files.
#
# Tool: K8S Secret Generator for HSM by H0llyW00dzZ
# Description: Creates Kubernetes secrets from a file, encoding values with base64 
#              (Note: Base64 automated by kubectl, when your k8s had HSM it will encrypted as well).
# Important:  Run this script where kubectl is installed. 
#             This script Required HSM and the format must be like "key=value", then run it where kubectl installed.
#
# Known Bugs (Bash bug): The following value format will cause issues: "key=value-value-value:value%!@value(value:value)/value?tls=value" it won't work       
#                        also note that it cannot be fixed with regular expressions due it bash problem (even it's possible, it just too complex), It might work-wells in unix-shellz.

# --- Configuration ---
ENV_FILE="worker-secret.txt" # Your secret file in "key=value" format (e.g, .env).
SECRET_NAME="my-worker-secrets"  # Name for the Kubernetes secret

# --- Script Logic ---

# Create the base kubectl command
kubectl_cmd="kubectl create secret generic $SECRET_NAME"

# Process the secret file
while IFS='=' read -r key value; do
  if [[ -z "$key" ]]; then continue; fi # Skip empty lines

  encoded_value=$(echo "$value") 
  kubectl_cmd="$kubectl_cmd --from-literal=$key=$encoded_value"
done < "$ENV_FILE"

# Execute the kubectl command
eval "$kubectl_cmd"

# --- Usage Example ---
# 1. Ensure your 'worker-secret.txt' file exists with content like:
#     DB_PASSWORD=your_db_password
#     API_KEY=your_api_key
# 2. Run the script: ./create_k8s_secret.sh 
