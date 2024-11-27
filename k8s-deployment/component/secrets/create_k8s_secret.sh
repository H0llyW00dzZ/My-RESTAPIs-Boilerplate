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

echo "$(tput setaf 4)
    _  __     _                          _            
   | |/ /    | |                        | |           
   | ' /_   _| |__   ___ _ __ _ __   ___| |_ ___  ___ 
   |  <| | | | '_ \ / _ \ '__| '_ \ / _ \ __/ _ \/ __|
   | . \ |_| | |_) |  __/ |  | | | |  __/ ||  __/\__ \\
   |_|\_\__,_|_.__/ \___|_|  |_| |_|\___|\__\___||___/
                                       
              $(tput sgr0) A Secret Tools Generator by H0llyW00dz

$(tput setaf 3)  Note: The format must be like 'key=value' in the secret file (e.g., .env). Run it where kubectl is installed.
        When your Kubernetes cluster has an HSM, it will be encrypted as well. Also note that regarding capabilities,
        it can handle up to 1K+++ more 'key=value' lines, depending on the requirements, instead of doing it manually.
        So, get good at bash scripting.$(tput sgr0)
"

# --- Configuration ---
read -p "Enter the secret file name: " ENV_FILE
read -p "Enter the Kubernetes secret name: " SECRET_NAME
read -p "Enter the Kubernetes namespace (default: default): " NAMESPACE

# Set default namespace if not provided
if [ -z "$NAMESPACE" ]; then
  NAMESPACE="default"
fi

# --- Script Logic ---

# Check if the secret file exists
if [ ! -f "$ENV_FILE" ]; then
  echo "$(tput setaf 1)Secret file '$ENV_FILE' does not exist. Please provide a valid file.$(tput sgr0)"
  exit 1
fi

# Create the base kubectl command
kubectl_cmd="kubectl create secret generic $SECRET_NAME -n $NAMESPACE"

# Process the secret file
while IFS='=' read -r key value; do
  if [[ -z "$key" ]]; then continue; fi # Skip empty lines

  encoded_value=$(echo "$value") 
  kubectl_cmd="$kubectl_cmd --from-literal=$key=$encoded_value"
done < "$ENV_FILE"

# Execute the kubectl command
eval "$kubectl_cmd"

# --- Usage Example ---
# 1. Run the script: ./create_k8s_secret.sh
# 2. Enter the secret file name when prompted (e.g., worker-secret.txt)
# 3. Enter the desired Kubernetes secret name when prompted (e.g., my-worker-secrets)
# 4. Enter the Kubernetes namespace when prompted (e.g., my-namespace) or leave it empty to use the default namespace

# --- Supported/Compatible ---
# This script should be compatible with any Bash/Shell environment on any operating system. As my primary personal use, I am using Git Bash on Windows.
