# K8s Deployment for REST API Boilerplate

<p align="center">
   <img src="https://kubernetes.io/images/kubernetes.png" alt="sailing-with-k8s" width="80">
</p>

This repository provides Kubernetes deployment files for a REST API boilerplate application. The primary focus is on leveraging Kubernetes to address critical infrastructure issues, enabling seamless scaling, and offering freedom from overpaying for licensing or other drama (e.g., bad competitors). Kubernetes allows you to pay primarily for direct hardware resources such as CPU, RAM, and disk, giving you the flexibility and efficiency needed for modern applications.

By utilizing Kubernetes, this deployment can handle billions of requests efficiently through Horizontal Pod Autoscaling (HPA), making it suitable for large-scale applications. The architecture is designed to be stateless, promoting scalability and stability across multiple nodes.

> [!NOTE]
> Without Kubernetes, this boilerplate cannot effectively address critical infrastructure issues such as scaling, security, and other (e.g., experimental solutions). It is designed to be stateless because, in Kubernetes, you can separate components instead of combining everything into a single stateful entity. This allows for seamless integration of components like databases and storage.

## List Documentation

### Redis Insight

This section covers the deployment of Redis Insight.

- [Documentation for Redis Insight](REDIS.md)

### REST APIs (This Repository)

This section covers the deployment of the REST API boilerplate application.

- [Documentation for REST APIs](RESTAPIs.md)

### MySQL

This section covers the deployment of the MySQL database service.

- [Documentation for MySQL](MySQL.md)

---

> [!NOTE]
> For users in Indonesia (🇮🇩), it's recommended to use the Singapore region instead of Indonesia when creating a cluster with a cloud provider. This is due to the traditional instability of the network in Indonesia (our home), which can lead to issues such as latency and packet loss.
