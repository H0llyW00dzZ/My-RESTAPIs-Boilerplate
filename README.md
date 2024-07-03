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
- Boring TLS 1.3 Protocol (WIP)
- Certificate Transparency (CT) Log

> [!NOTE]  
> `Boring TLS 1.3 Protocol` This project utilizes the industry-standard TLS 1.3 protocol, known for its exceptional security and reliability. We call it `Boring TLS` because its robust security makes it almost boring.

> [!NOTE]
> `Certificate Transparency (CT) Log` This project includes built-in support for submitting certificates to Certificate Transparency logs. It supports certificates with both RSA and ECDSA keys. By submitting certificates to CT logs, it enhances the security and transparency of the TLS ecosystem, allowing domain owners to monitor and detect fraudulent certificates or other bad certificates. The CT log functionality is implemented using the `SubmitToCTLog` function, which takes the certificate, private key, and CT log details as input. It encodes the certificate, calculates the hash, creates a JSON payload, sends an HTTP request to the CT log server, and verifies the signed certificate timestamp (SCT) received in the response. The `VerifySCT` method is responsible for verifying the SCT based on the public key type (RSA or ECDSA) and ensures the integrity of the timestamp. With this feature, it is easy to integrate certificate transparency into the TLS setup and contribute to a more secure web.

> [!WARNING]
> Some features might not work as expected when running this repo on certain cloud providers. For example, the `Rich Presence TUIs` feature requires a [`tty`](https://en.wikipedia.org/wiki/Tty_(Unix)).
> The reason why some features might not work as expected after implementation is because this repo is designed to be `top-level` and follow the `Unix philosophy`. Plus, most of the `Go code` in this repo follows many `best practices` and idioms from `Effective Go`.

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

<img src="https://i.imgur.com/FzESWLM.png" alt="tls-go1.22.5">
<img src="https://i.imgur.com/fGBOy2k.png" alt="tls-go1.22.5">


## Architecture

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

## License

This project is dual-licensed under the BSD 3-Clause License and the MIT License - see the [LICENSE](LICENSE) file for details.
