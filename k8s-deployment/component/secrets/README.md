# Kubernetes Secret Generator (HSM-Ready)

<p align="center">
   <img src="https://kubernetes.io/images/kubernetes.png" alt="sailing-with-k8s" width="80">
   <img src="https://raw.githubusercontent.com/kubernetes/community/refs/heads/master/icons/png/resources/labeled/secret-128.png" alt="The-Black-Pearl" width="80">
</p>

A simple Bash script to generate Kubernetes secrets from a file containing `key=value` pairs. This tool automates the creation of Kubernetes secrets, encoding values with `base64` (automatically handled by `kubectl`) and ensuring compatibility with Kubernetes clusters that support Hardware Security Modules (HSM).

---

## Features

- Converts a file of `key=value` pairs (e.g., `.env` file) into a Kubernetes secret.
- Supports Kubernetes clusters with HSM for additional encryption.
- Handles large files with thousands of `key=value` pairs.
- Allows specifying the Kubernetes namespace and secret name.
- Simplifies secret management by automating the process.

---

## Requirements

- **Bash**: Ensure you are running this script in a Bash-compatible shell.
- **kubectl**: The Kubernetes CLI tool must be installed and configured to connect to your cluster.
- **HSM Support**: If your Kubernetes cluster has HSM enabled, the secrets will be encrypted automatically.

---

## Installation

1. Make the script executable:
   ```bash
   chmod +x create_k8s_secret.sh
   ```

2. (Optional) Install `dos2unix` to ensure the script has proper Unix line endings:
   ```bash
   sudo apt update
   sudo apt install dos2unix
   dos2unix create_k8s_secret.sh
   ```

---

## Usage

1. Run the script:
   ```bash
   ./create_k8s_secret.sh
   ```

2. Follow the prompts:
   - **Enter the secret file name**: Provide the path to the file containing `key=value` pairs (e.g., `.env` file).
   - **Enter the Kubernetes secret name**: Specify the name of the Kubernetes secret you want to create.
   - **Enter the Kubernetes namespace**: Provide the namespace where the secret should be created (default is `default`).

3. The script will process the file and create the Kubernetes secret.

---

## Input File Format

The input file must be formatted as `key=value` pairs, one per line. For example:

```env
DATABASE_URL=postgres://user:password@host:5432/dbname
API_KEY=your-api-key
SECRET_TOKEN=your-secret-token
```

---

## Example

1. Create a `.env` file:
   ```env
   DB_USER=admin
   DB_PASSWORD=supersecret
   API_TOKEN=abcdef123456
   ```

2. Run the script:
   ```bash
   ./create_k8s_secret.sh
   ```

3. Enter the following when prompted:
   - Secret file name: `./.env`
   - Kubernetes secret name: `my-app-secrets`
   - Kubernetes namespace: `my-namespace`

4. The script will generate a Kubernetes secret:
   ```bash
   kubectl create secret generic my-app-secrets -n my-namespace \
       --from-literal=DB_USER=admin \
       --from-literal=DB_PASSWORD=supersecret \
       --from-literal=API_TOKEN=abcdef123456
   ```

---

## Known Bugs and Limitations

- **Unsupported Value Formats**: Certain complex value formats (e.g., `key=value-value:value%!@value`) may cause issues due to Bash limitations.
- **Empty Lines**: Empty lines in the input file are skipped.
- **Windows Line Endings**: If the script was created or edited on Windows, ensure it has Unix-style line endings using `dos2unix`.

---

## Troubleshooting

### Error: `/bin/bash: ‘bash\r’: No such file or directory`
This occurs when the script has Windows-style line endings (`\r\n`) instead of Unix-style (`\n`). Fix it with:
```bash
dos2unix create_k8s_secret.sh
```

### Error: `kubectl: command not found`
Ensure `kubectl` is installed and configured:
```bash
sudo apt install -y kubectl
kubectl config view
```
