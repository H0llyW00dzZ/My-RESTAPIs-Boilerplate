# My RESTful API Boilerplate

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

## TODO Features

- [ ] Custom Monitoring Integration with Kubernetes ecosystems, such as Grafana, etc. at the top level

- [X] Translate to any language

> [!NOTE]  
> Additional features will be added later as my primary focus is on using this with the Fiber framework.


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

- `MySQL` without `Redis` + 100K `POST` Request

<img src="https://i.imgur.com/C9hZYDz.png" alt="Memory Usage with MYSQL">

> [!NOTE]
> The `Resource Memory Usage` section demonstrates how `Go` has stable and low memory overhead compared to other languages, especially `Java` (See [this article](https://learn.microsoft.com/en-us/azure/spring-apps/enterprise/concepts-for-java-memory-management) for more information on Java memory management Lmao.) hahaha.

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
