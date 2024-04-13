# My RESTful API Boilerplate

This project provides a boilerplate for building RESTful APIs using Go. It is designed to a quick start with best practices, easy configuration, and a clear project structure.

## Features

- Use Fiber Framework
- Pre-configured MySQL and Redis integration
- Middleware examples for logging and authentication
- Scalable project structure
- Dual-licensed under BSD 3-Clause and MIT (specific files)
- High Performance (Thanks to the [`Sonic JSON serializing & deserializing library`](https://github.com/bytedance/sonic))

## TODO Features

- [ ] Custom Monitoring Integration with Kubernetes ecosystems, such as Grafana, etc. at the top level


> [!NOTE]  
> Additional features will be added later as my primary focus is on using this with the Fiber framework.

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
> The `Resource Memory Usage` section demonstrates how `Go` has stable and low memory overhead compared to other languages, especially `Java` hahaha.

## License

This project is dual-licensed under the BSD 3-Clause License and the MIT License - see the [LICENSE](LICENSE) file for details.
