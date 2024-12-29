# My RESTful API Boilerplate

<p align="center">
  <img src="https://i.imgur.com/Rm7I8uK.png" alt="Golang logo" width="500">
  <br>
  <i>Image Copyright © <a href="https://github.com/SAWARATSUKI">SAWARATSUKI</a>. All rights reserved.</i>
  <br>
  <i>Image used under permission from the copyright holder.</i>
</p>

This project provides a boilerplate for building RESTful APIs using Go. It is designed to a quick start with best practices, easy configuration, and a clear project structure.

> [!WARNING]  
> This boilerplate is specifically tailored for extensive codebases and scalable RESTful API applications, along with frontend websites that may span multiple sites.
> Its design ensures scalability and straightforward maintenance, making it highly suitable for complex projects. However, it may not be the ideal choice for smaller-scale RESTful API applications or single-website frontends where such robust architecture is not required.

> [!TIP]
> Given that this boilerplate focuses on `high performance` (e.g., `it can handle substantial workloads`), it is recommended to `run and isolate` it on [`Kubernetes`](https://kubernetes.io/). It has been tested alongside the [`worker package`](https://github.com/H0llyW00dzZ/My-RESTAPIs-Boilerplate/tree/master/worker) and has proven to be `fully stable` for `low to high workloads (concurrency)`. `Running and isolating` it on [`Kubernetes`](https://kubernetes.io/) allows for easy scaling when additional `resources`, such as `vCPUs`, are needed.

## Features

- Use Fiber Framework
- Pre-configured MySQL and Redis integration
- Middleware examples for logging and authentication
- Scalable project structure
- Dual-licensed under BSD 3-Clause and MIT (specific files)
- High Performance (Thanks to the [`Sonic JSON serializing & deserializing library`](https://github.com/bytedance/sonic))
- Cryptography Techniques
- High Quality Go Codes
- Rich Presence TUIs using [`charm.sh`](https://charm.sh/)
- ~~Boring TLS 1.3 Protocol (WIP)~~ -> EOL (No longer actively developed due to time constraints)
- ~~Certificate Transparency (CT) Log~~ -> EOL (No longer actively developed due to time constraints)

> [!NOTE]  
> `Boring TLS 1.3 Protocol` This project utilizes the industry-standard TLS 1.3 protocol, known for its exceptional security and reliability. We call it `Boring TLS` because its robust security makes it almost boring.

> [!NOTE]
> `Certificate Transparency (CT) Log` This project includes built-in support for submitting certificates to Certificate Transparency logs. It supports certificates with both RSA and ECDSA keys. By submitting certificates to CT logs, it enhances the security and transparency of the TLS ecosystem, allowing domain owners to monitor and detect fraudulent certificates or other bad certificates. The CT log functionality is implemented using the `SubmitToCTLog` function, which takes the certificate, private key, and CT log details as input. It encodes the certificate, calculates the hash, creates a JSON payload, sends an HTTP request to the CT log server, and verifies the signed certificate timestamp (SCT) received in the response. The `VerifySCT` method is responsible for verifying the SCT based on the public key type (RSA or ECDSA) and ensures the integrity of the timestamp. With this feature, it is easy to integrate certificate transparency into the TLS setup and contribute to a more secure web.

> [!WARNING]
> Some features might not work as expected when running this repo on certain cloud providers. For example, the `Rich Presence TUIs` feature requires a [`tty`](https://en.wikipedia.org/wiki/Tty_(Unix)).
> The reason why some features might not work as expected after implementation is because this repo is designed to be `top-level` and follow the `Unix philosophy`. Plus, most of the `Go code` in this repo follows many `best practices` and idioms from `Effective Go`.

> [!TIP]
> Due to the large codebase of this repository, which includes many sub-packages such as helper functions (lazy splitting into another repository), consider implementing your own or an alternative `main` Go mechanism ([`follow best practices here`](https://go.dev/doc/modules/layout#server-project)). This can be used across multiple containers in a single deployment for infrastructure purposes. Personally, I have created many such mechanisms for running in [`Kubernetes (k8s)`](https://kubernetes.io/) based on this starter repository.

## TODO Features

Move [`here`](https://github.com/users/H0llyW00dzZ/projects/2/views/3?sliceBy%5Bvalue%5D=Todo)

> [!NOTE]  
> Additional features will be added later as my primary focus is on using this with the `Fiber framework` and `Cryptography Techniques`.


## Supported Additional Features

> [!NOTE]  
> The following list includes additional features that I've been using before for extremely scalable and high-performance applications:

- [`Protocol Buffers`](https://protobuf.dev/)
- [`Templ`](https://templ.guide/)

## Motivation

The motivation for sharing this RESTful API Boilerplate on GitHub is to streamline my development process, ensuring I don't have to start from scratch each time. It's a robust, secure foundation that adheres to best practices, designed for quick and efficient project initiation.

## Philosophy

This boilerplate is grounded in the Unix philosophy, emphasizing simplicity, modularity, and reusability. Each component is designed to do one thing and do it well, with the ability to combine components seamlessly to perform complex tasks. This approach ensures that the boilerplate remains lightweight, efficient, and easy to maintain.

<p align="center">
  <img src="https://i.imgur.com/jIX8esm.png" alt="Gopher Drink" width="250">
</p>

## Resource Memory Usage

- `MySQL` without `Redis` + 100K `POST` Request (`Outdated` - `go1.22.2`)

<img src="https://i.imgur.com/C9hZYDz.png" alt="Memory Usage with MYSQL">

> [!NOTE]
> The `Resource Memory Usage` section demonstrates how `Go` has stable and low memory overhead compared to other languages, especially `Java` (See [this article](https://learn.microsoft.com/en-us/azure/spring-apps/enterprise/concepts-for-java-memory-management) for more information on Java memory management Lmao.) hahaha.

- `Idle` (`Outdated` - `go1.22.2`)

<img src="https://i.imgur.com/1dG4E9G.png" alt="Idle">

> [!NOTE]
> The `Idle` section demonstrates the memory usage when there is no request. The maximum average memory usage is around 21.5MB. Go routines (100 goroutines) along with a semaphore are used to automatically handle high traffic situations (e.g., lots of requests).
>
> Also Note that even with high incoming traffic (e.g., `1 million requests`), the maximum average memory usage is still relatively `low`, around `100MB` (e.g., `50MB`), due to the use of the `semaphore`.

### ***Latest***

- `Idle` (`Outdated` - `go1.22.3`)

<img src="https://i.imgur.com/9AstMIl.png" alt="Idle-go1.22.3">

```sh
sample#memory_total=8.01MB
sample#memory_rss=7.92MB
sample#memory_cache=0.09MB
sample#memory_swap=0.00MB
sample#memory_pgpgin=0pages
sample#memory_pgpgout=0pages
sample#memory_quota=1024.00MB
```

> [!NOTE]
> **Go Memory Usage Improvements:**
> - The `Idle` memory usage in the latest version of `Go` has been optimized, reducing it to approximately `7.92MB` compared to the previous version.
>
> **Average Memory Usage Breakdown:**
> - The average memory usage of `11.1MB` includes:
>   - Go routine scheduler for database operations (separate from the 100 goroutines used for handling high incoming traffic automatically with a semaphore)
>   - Load for the `front-end` website
>   - Load for the `REST APIs` (`Full Stack`)
> - The application remains stable without any performance issues, such as bottlenecks or memory leaks, and has zero vulnerabilities.
>
> **Go Performance Comparison:**
> - Compared to other languages, particularly for full-stack development, `Go` demonstrates superior performance.
>
> <p align="center">
>   <img src="https://i.imgur.com/PxjZ0Dz.png" alt="gopher run" />
> </p>

#### ***TLS 1.3 Connection Stable***

- `Standard TLS 1.3` (`Latest Go Version` - `go1.22.5`)

<img src="https://i.imgur.com/NiUdm8I.png" alt="tls-go1.22.5">
<img src="https://i.imgur.com/jHgqAey.png" alt="tls-go1.22.5">
<img src="https://i.imgur.com/OUpBdUm.png" alt="tls-go1.22.5-cracked-zer0ms-response">
<img src="https://i.imgur.com/yx6HBU8.png" alt="tls-go1.22.5-cracked-zer0ms-response">

> [!NOTE]
> The screenshots provided demonstrate a `Standard TLS 1.3` connection using the `Closed networks or intranets` method, which involves Cloudflare (`Frontend`) and Heroku (`Backend`). Without Cloudflare, any browser or tools like curl won't be able to access the backend server.


## Frontend Performance

- `Error Page` (`Include Wildcard Handler`)

<img src="https://i.imgur.com/gbGM6bS.png" alt="error-page">

> [!NOTE]
> The screenshot provided demonstrates the performance of an `Error Page`, and it can easily be optimized to achieve all green metrics.

> [!TIP]
> Also note that when optimized to achieve all green metrics which is easily, especially in the SEO category, it can be beneficial for business logic purposes. (e.g, Search engines like `Google` and `Microsoft Bing` tend to favor websites with good performance metrics, which can lead to improved search rankings and increased visibility.)


## Prometheus

- `Latest Go Version`

<img src="https://i.imgur.com/KeHzVn0.png" alt="prometheus-portable">
<img src="https://i.imgur.com/U6o6x9U.png" alt="prometheus-portable">
<img src="https://i.imgur.com/UySPjRn.png" alt="prometheus-portable">
<img src="https://i.imgur.com/azGS0et.png" alt="prometheus-portable">

> [!NOTE]
> The `grafana` dashboard is not shareable due to it being bound to my security configurations for real-world monitoring in other production environments.

## Architecture (Tree)

Below is the architecture of this boilerplate and how it looks. I created this for REST APIs about volcano 🌋 monitoring used by the government (has been done before), so it can easily monitor volcanoes in the real world.

- Backend (Pure `Go`)

```
backend/
|-- cmd/
|   `-- server/ (main application entry point)
|       `-- main.go
|-- pkg/
|   |-- restapis/ (API route handlers)
|   |-- any/ (any related code)
|-- internal/ (private application and library code)
|   |-- any/ (any related code)
|-- .env (optional environment variables since it can placed anywhere)
`-- go.mod (dependencies)
```

- Frontend (Any Framework `TS` e.g, React,NextJS or pure `JS` for `HTMX`)

#### Example:

```
frontend/
|-- pages/
|   |-- api/ (optional, for Next.js API routes if needed)
|   |-- _app.js (global page layouts and state)
|   |-- index.js (home page)
|   `-- [...other pages]
|-- public/ (static files like images, fonts)
|-- src/
|   |-- components/ (shared React components)
|   |   `-- [various components]
|   |-- styles/ (global styles, theme)
|   |-- hooks/ (custom React hooks)
|   |-- utils/ (utility functions)
|   |-- lib/ (libraries and configurations)
|   `-- context/ (React context files for state management)
|-- .env.local (environment variables)
|-- next.config.js (Next.js configuration)
`-- package.json (dependencies and scripts)
```

#### Current Boilerplate Tree:

```
├── Dockerfile
├── LICENSE
├── README.md
├── SECURITY.md
├── backend
│   ├── cmd
│   │   └── server
│   │       ├── run.go
│   │       ├── run_heroku.go
│   │       ├── run_immutable.go
│   │       └── run_tls_fips.go
│   ├── internal
│   │   ├── database
│   │   │   ├── auth.go
│   │   │   ├── backup.go
│   │   │   ├── cloudflare-kv
│   │   │   │   └── setup.go
│   │   │   ├── constant.go
│   │   │   ├── helper.go
│   │   │   ├── mysql_redis.go
│   │   │   ├── redis_json.go
│   │   │   ├── setup.go
│   │   │   ├── sql_injection_test.go
│   │   │   └── tls.go
│   │   ├── logger
│   │   │   ├── constant.go
│   │   │   └── logger.go
│   │   ├── middleware
│   │   │   ├── authentication
│   │   │   │   ├── crypto
│   │   │   │   │   ├── bcrypt
│   │   │   │   │   │   ├── bcrypt.go
│   │   │   │   │   │   ├── bcrypt_test.go
│   │   │   │   │   │   ├── compare_password.go
│   │   │   │   │   │   ├── docs.go
│   │   │   │   │   │   └── hash_password.go
│   │   │   │   │   ├── cipher.go
│   │   │   │   │   ├── crypto.go
│   │   │   │   │   ├── crypto_test.go
│   │   │   │   │   ├── deadcode.go
│   │   │   │   │   ├── decrypt.go
│   │   │   │   │   ├── encrypt.go
│   │   │   │   │   ├── gopherpocket
│   │   │   │   │   │   └── keyrotation
│   │   │   │   │   │       ├── docs.go
│   │   │   │   │   │       ├── gopherkey.go
│   │   │   │   │   │       └── gopherkey_test.go
│   │   │   │   │   ├── gpg
│   │   │   │   │   │   ├── benchmark_test.go
│   │   │   │   │   │   ├── config.go
│   │   │   │   │   │   ├── docs.go
│   │   │   │   │   │   ├── encrypt.go
│   │   │   │   │   │   ├── encrypt_test.go
│   │   │   │   │   │   ├── file.go
│   │   │   │   │   │   ├── key_info.go
│   │   │   │   │   │   ├── key_ring.go
│   │   │   │   │   │   ├── keybox.go
│   │   │   │   │   │   ├── keybox_test.go
│   │   │   │   │   │   └── new.go
│   │   │   │   │   ├── helper.go
│   │   │   │   │   ├── hybrid
│   │   │   │   │   │   ├── decryptcookie.go
│   │   │   │   │   │   ├── encryptcookie.go
│   │   │   │   │   │   ├── hybrid.go
│   │   │   │   │   │   ├── hybrid_stream.go
│   │   │   │   │   │   ├── hybrid_test.go
│   │   │   │   │   │   └── stream
│   │   │   │   │   │       ├── benchmark_test.go
│   │   │   │   │   │       ├── chunk.go
│   │   │   │   │   │       ├── decrypt_stream.go
│   │   │   │   │   │       ├── digest.go
│   │   │   │   │   │       ├── docs.go
│   │   │   │   │   │       ├── encrypt_stream.go
│   │   │   │   │   │       ├── new_stream.go
│   │   │   │   │   │       ├── nonce.go
│   │   │   │   │   │       └── stream_test.go
│   │   │   │   │   ├── keyidentifier
│   │   │   │   │   │   ├── docs.go
│   │   │   │   │   │   ├── ecdsa_sign.go
│   │   │   │   │   │   ├── generate.go
│   │   │   │   │   │   ├── keyidentifier_test.go
│   │   │   │   │   │   └── new.go
│   │   │   │   │   ├── rand
│   │   │   │   │   │   ├── benchmark_test.go
│   │   │   │   │   │   ├── fixed_size.go
│   │   │   │   │   │   ├── rand_test.go
│   │   │   │   │   │   └── uuid.go
│   │   │   │   │   ├── tls
│   │   │   │   │   │   ├── docs.go
│   │   │   │   │   │   └── setup.go
│   │   │   │   │   ├── vault
│   │   │   │   │   │   ├── new.go
│   │   │   │   │   │   └── transit.go
│   │   │   │   │   ├── web3
│   │   │   │   │   │   └── eth
│   │   │   │   │   │       ├── docs.go
│   │   │   │   │   │       └── new.go
│   │   │   │   │   └── webauthn
│   │   │   │   │       ├── docs.go
│   │   │   │   │       ├── login.go
│   │   │   │   │       ├── protocol.go
│   │   │   │   │       └── registration.go
│   │   │   │   ├── helper
│   │   │   │   │   ├── constant.go
│   │   │   │   │   └── keyauth.go
│   │   │   │   └── keyauth
│   │   │   │       ├── constant.go
│   │   │   │       ├── error.go
│   │   │   │       ├── success.go
│   │   │   │       └── validator.go
│   │   │   ├── constant.go
│   │   │   ├── csp
│   │   │   │   ├── config.go
│   │   │   │   ├── csp_test.go
│   │   │   │   └── new.go
│   │   │   ├── custom_logger_tag.go
│   │   │   ├── custom_next.go
│   │   │   ├── filesystem
│   │   │   │   └── crypto
│   │   │   │       └── signature
│   │   │   │           ├── hmac_sign.go
│   │   │   │           ├── hmac_test.go
│   │   │   │           └── hmac_verify.go
│   │   │   ├── frontend_routes.go
│   │   │   ├── helper.go
│   │   │   ├── init.go
│   │   │   ├── lb_test.go
│   │   │   ├── monitor
│   │   │   │   ├── docs.go
│   │   │   │   └── prometheus.go
│   │   │   ├── restapis_routes.go
│   │   │   ├── restime
│   │   │   │   ├── config.go
│   │   │   │   └── new.go
│   │   │   ├── router
│   │   │   │   ├── domain
│   │   │   │   │   ├── config.go
│   │   │   │   │   └── new.go
│   │   │   │   ├── logs
│   │   │   │   │   └── connection
│   │   │   │   │       ├── config.go
│   │   │   │   │       ├── conn_test.go
│   │   │   │   │       ├── get.go
│   │   │   │   │       └── new.go
│   │   │   │   └── proxytrust
│   │   │   │       ├── config.go
│   │   │   │       ├── new.go
│   │   │   │       └── proxytrust_test.go
│   │   │   ├── routes.go
│   │   │   ├── routes_immutable.go
│   │   │   ├── routes_non_immutable.go
│   │   │   ├── storage.go
│   │   │   └── utils.go
│   │   ├── server
│   │   │   ├── boringtls.go
│   │   │   ├── boringtls_cert.go
│   │   │   ├── boringtls_cert_test.go
│   │   │   ├── boringtls_test.go
│   │   │   ├── constant.go
│   │   │   ├── ct_verifier.go
│   │   │   ├── ct_verifier_test.go
│   │   │   ├── helper.go
│   │   │   ├── helper_tls_test.go
│   │   │   ├── k8s
│   │   │   │   ├── docs.go
│   │   │   │   └── metrics
│   │   │   │       └── prometheus.go
│   │   │   ├── mount_routes.go
│   │   │   ├── register_routes.go
│   │   │   └── startup_async.go
│   │   └── translate
│   │       └── language.go
│   └── pkg
│       ├── archive
│       │   ├── archive.go
│       │   ├── config.go
│       │   ├── do.go
│       │   ├── do_test.go
│       │   └── docs.go
│       ├── convert
│       │   ├── docs.go
│       │   ├── helper.go
│       │   ├── to_bytes.go
│       │   └── to_bytes_test.go
│       ├── gc
│       │   ├── docs.go
│       │   ├── reduce_http_client_overhead.go
│       │   ├── reduce_overhead.go
│       │   └── unique.go
│       ├── header
│       │   └── htmx
│       │       ├── constant.go
│       │       └── docs.go
│       ├── mime
│       │   ├── ascii.go
│       │   ├── docs.go
│       │   └── mime.go
│       ├── network
│       │   └── cidr
│       │       ├── validate.go
│       │       └── validate_test.go
│       └── restapis
│           ├── helper
│           │   ├── auth
│           │   │   ├── apikey.go
│           │   │   └── constant.go
│           │   ├── generate_apikey.go
│           │   ├── generate_apikey_test.go
│           │   ├── json
│           │   │   └── sonic
│           │   │       ├── config.go
│           │   │       └── docs.go
│           │   ├── numeric.go
│           │   ├── numeric_test.go
│           │   ├── restapis_error.go
│           │   └── restapis_error_test.go
│           └── server
│               └── health
│                   ├── cache.go
│                   ├── constant.go
│                   ├── db.go
│                   ├── helper.go
│                   ├── mysql.go
│                   └── redis.go
├── env
│   ├── docs.go
│   ├── env.go
│   └── getenv.go
├── frontend
│   ├── htmx
│   │   ├── error_page_handler
│   │   │   ├── 400.templ
│   │   │   ├── 400_templ.go
│   │   │   ├── 401.templ
│   │   │   ├── 401_templ.go
│   │   │   ├── 403.templ
│   │   │   ├── 403_templ.go
│   │   │   ├── 404.templ
│   │   │   ├── 404_templ.go
│   │   │   ├── 500.templ
│   │   │   ├── 500_templ.go
│   │   │   ├── 502.templ
│   │   │   ├── 502_templ.go
│   │   │   ├── 503.templ
│   │   │   ├── 503_templ.go
│   │   │   ├── 504.templ
│   │   │   ├── 504_templ.go
│   │   │   ├── base.templ
│   │   │   ├── base_templ.go
│   │   │   ├── page_handler.go
│   │   │   ├── render_frontend.go
│   │   │   └── static_handler.go
│   │   ├── public
│   │   │   └── assets
│   │   │       └── css
│   │   │           └── base-tailwind.css
│   │   └── site
│   │       ├── footer.templ
│   │       ├── footer_templ.go
│   │       ├── head.templ
│   │       ├── head_templ.go
│   │       ├── header.templ
│   │       ├── header_templ.go
│   │       ├── script.templ
│   │       └── script_templ.go
│   └── public
│       ├── assets
│       │   ├── css
│       │   │   ├── base-tailwind.css
│       │   │   └── raw.css
│       │   ├── images
│       │   │   ├── android-chrome-192x192.png
│       │   │   ├── android-chrome-512x512.png
│       │   │   ├── apple-touch-icon.png
│       │   │   ├── browserconfig.xml
│       │   │   ├── favicon-16x16.png
│       │   │   ├── favicon-32x32.png
│       │   │   ├── favicon.ico
│       │   │   ├── http_error_codes
│       │   │   │   ├── 403-Forbidden.png
│       │   │   │   ├── 404-NotFound.png
│       │   │   │   ├── 500-InternalServerError.png
│       │   │   │   └── 503-ServiceUnavailable.png
│       │   │   ├── logo
│       │   │   │   └── gopher-run.png
│       │   │   ├── mstile-150x150.png
│       │   │   ├── safari-pinned-tab.svg
│       │   │   └── site.webmanifest
│       │   └── js
│       │       ├── htmx.indicator.min.js
│       │       ├── htmx.min.js
│       │       └── tailwind.min.dark.js
│       └── magic_embedded.go
├── gcloud-builds.yaml
├── go.mod
├── go.sum
├── k8s-deployment
│   ├── MySQL.md
│   ├── README.md
│   ├── REDIS.md
│   ├── RESTAPIs.md
│   ├── component
│   │   ├── coredns
│   │   │   ├── coredns-custom-resource.yaml
│   │   │   └── coredns-log.yaml
│   │   └── secrets
│   │       └── create_k8s_secret.sh
│   ├── mysql-deploy-cpu-boost.yaml
│   ├── mysql-deploy.yaml
│   ├── mysql-storage-eks.yaml
│   ├── mysql-storage-gke.yaml
│   ├── prometheus_portable.yaml
│   ├── prometheus_portable_rules_record.yaml
│   ├── redis-insight.yaml
│   ├── restapis-deploy.yaml
│   ├── restapis-ingress.yaml
│   └── traffic
│       └── nginx
│           ├── ingress-nginx-configmap.yaml
│           ├── ingress-nginx-hpa.yaml
│           ├── ingress-nginx-priority.yaml
│           └── ingress-nginx-vpa.yaml
├── tailwind.config.js
├── translate.json
└── worker
    ├── config.go
    ├── do_work.go
    ├── docs.go
    ├── jobs.go
    ├── token_bucket.go
    └── worker_test.go

79 directories, 253 files
```

> [!NOTE]
> The `Current Boilerplate Tree` is designed as a `modular framework`, and it is easily maintainable even if it reaches 1K files. Personally, I've been maintaining over 1k files as well
> from this boilerplate, and it runs smoothly on Kubernetes (K8s) ⛵ ☸.

## Git Mirror (Auto Synced)

### Server Location: Singapore (Stable and very low latency for `Southeast Asian Regions`, ranging from `10ms ~ 50ms`)

- ![gitea](https://git.b0zal.io/assets/img/logo.svg) [My REST APIs Boilerplate](https://git.b0zal.io/H0llyW00dzZ/My-RESTAPIs-Boilerplate.git)

#### 💻 >_ Shell (git protocol https):

```sh
git clone https://git.b0zal.io/H0llyW00dzZ/My-RESTAPIs-Boilerplate.git
```

> [!NOTE]
> <p align="center">
>   <img src="https://kubernetes.io/images/kubernetes.png" alt="fully-managed-and-isolated-by-k8s" width="80">
> </p>
>
> **Server Management and Isolation by Kubernetes:**
> - The `Storage` is secure and fully encrypted (`end-to-end`), designed with flexibility in mind. It is suitable for `automated attach/detach` processes within a `cluster`.
> - The `Network` utilizes a `network load balancer` controlled by `Ingress Nginx`, optimizing latency for the `APAC` region, ensuring smooth sailing ⛵ ☸.
> - The `Git Protocol SSH` should function properly as it utilizes the `TCP Service Nginx`. It is `fully secured`, making it resistant to `brute-force attacks` and `exploits` due to its `underlying logic`, which incorporates [`Elliptic Curve Cryptography (ECC)`](https://en.wikipedia.org/wiki/Elliptic-curve_cryptography).
>
> Note that The `Git Protocol SSH` is currently disabled because the new load balancer is not functioning properly with the TCP service in ingress-nginx (not security issues related CVE vulnerabilities). Rest assured, the git mirror is secure and respects privacy (100%), as it runs on my cluster (very secure) used for development, sandbox, and production. You can explore other repositories [here](https://git.b0zal.io/explore/repos).

> [!TIP]
> For those in `Indonesia`, if you are unable to clone repositories (e.g., using `git clone`) from the [`Git Mirror (Auto Synced)`](https://git.b0zal.io/H0llyW00dzZ/My-RESTAPIs-Boilerplate.git), try using a VPN. This issue may be related to the new load balancer, which could be blocked (by internet provider) or filtered by my firewall mechanism.
>
> It is also more faster with the new load balancer, even when using a VPN (`e.g., VPN in Singapore, Indonesia, Malaysia, and other APAC regions`), and this improvement applies to both `mobile devices and desktops`, which incorporate [`Elliptic Curve Cryptography (ECC)`](https://en.wikipedia.org/wiki/Elliptic-curve_cryptography) for `HTTPS/TLS`.

# Supported Architectures

Due to the use of [Sonic JSON](https://github.com/bytedance/sonic) for encoding/decoding, this boilerplate supports only the following architectures:
- `AMD64`
- `ARM64`

> [!NOTE]
> It can also be ideal for `multi-architecture` workloads.

## License

This project is dual-licensed under the BSD 3-Clause License and the MIT License - see the [LICENSE](LICENSE) file for details.
