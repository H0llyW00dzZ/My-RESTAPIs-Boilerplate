# Copyright (c) 2024 H0llyW00dz All rights reserved.
# 
# License: BSD 3-Clause License
#
# Tool: K8S Secret Hardware security module (HSM) Generator by H0llyW00dzZ
# Important: This Secret Generator for K8S Required HSM and the format must be like "key=value", then run it where kubectl installed.
#
# Example Format:
# key=value
# key2=value2
#
# Note: There it's no limit, so it's capable handle 1k Line
#
#!/bin/bash

# Your secret file (e.g, .env)
ENV_FILE="worker-secret.txt"

# Create a command string for kubectl
kubectl_cmd="kubectl create secret generic my-worker-secrets" # change this "my-worker-secrets" 

# Loop through the worker-secret.txt file, adding --from-literal arguments
while read line; do
  key="${line%%=*}"
  value="${line#*=}"
  encoded_value=$(echo "$value")
  kubectl_cmd="$kubectl_cmd --from-literal=$key=$encoded_value"
done < "$ENV_FILE"

# Execute the kubectl command
eval "$kubectl_cmd"
