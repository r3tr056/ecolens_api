# Ecoview API Contribution Guidelines

Thank you for considering contributing to Ecoview API! We appreciate your interest in making our project better.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
  - [Setting Up the Development Environment](#setting-up-the-development-environment)
  - [Forking the Repository](#forking-the-repository)
  - [Branching Strategy](#branching-strategy)
- [Development Guidelines](#development-guidelines)
  - [Code Style](#code-style)
  - [Testing](#testing)
  - [Documentation](#documentation)
- [Pull Requests](#pull-requests)
- [Issues](#issues)
- [Communication](#communication)
- [License](#license)

## Code of Conduct

Please note that this project is governed by our [Code of Conduct](link/to/code-of-conduct.md). Be sure to review and adhere to these guidelines.

## Getting Started

### Setting Up the Development Environment

Ecoview API is written in Golang and uses GoFiber as the web framework. It also relies on PostgreSQL and Redis for data storage and Pub/Sub for messaging in a microservices architecture.

1. Install [Golang](https://golang.org/doc/install).
2. Install [Docker](https://www.docker.com/get-started) for running PostgreSQL and Redis containers.
3. Fork and clone the Ecoview API repository.

### Forking the Repository

If you haven't already, fork the Ecoview API repository on GitHub. This will create your copy of the project.

### Branching Strategy

Create a new branch for each contribution. Use a descriptive branch name that reflects the purpose of your changes.

```bash
git checkout -b feature/new-feature
```

## Development Guidelines

### Code Style

Follow the Go [official coding conventions](https://golang.org/doc/effective_go.html) and the specific guidelines mentioned in the codebase. Use tools like `gofmt` and `golint` to maintain code consistency.

### Testing

Write comprehensive tests for new features and ensure that existing tests pass. Use tools like `go test` to run tests.

### Documentation

Maintain clear and concise code comments. If you introduce new features, update the project's documentation accordingly. For API changes, ensure that the OpenAPI specification is updated.

## Pull Requests

1. Ensure your code adheres to the guidelines mentioned above.
2. Test your changes thoroughly.
3. Update the project documentation if necessary.
4. Open a pull request against the `main` branch of the original Ecoview API repository.

## Issues

If you encounter any issues or have suggestions, feel free to open a GitHub issue. Provide detailed information about the problem and steps to reproduce it.

## Communication

Reach out our development team at [dev.ecoview@gmail.com](mailto:dev.ecoview@gmail.com) to discuss ideas and ask questions.
