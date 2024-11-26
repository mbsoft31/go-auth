# Auth Package Documentation

Welcome to the `auth` package documentation! This package provides a lightweight, SQLite-backed authentication and session management system designed for Go applications.

## Features
- User management: creation, retrieval, and updates
- Session management with cookies
- Middleware for Echo framework
- Configurable options for cookies and security
- Secure password hashing

## **Documentation Structure**
- [Auth](./auth.md): Overview of user and session management.
- [Middleware](./middleware.md): Explanation of the Echo middleware for session authentication.
- [Store](./store.md): Details of database interactions and available store methods.
- [Config](./config.md): Configuration settings and defaults for the package.

## **Getting Started**
1. Clone the repository and install dependencies.
2. Initialize the store with a valid configuration.
3. Integrate middleware for protected routes.
4. Refer to the individual documentation files for detailed usage.

## **Quick Links**
- [Setup](#getting-started)
- [API Reference](./docs)
- [Examples](./examples)
- [Contributing](./CONTRIBUTING.md)

## **License**
This package is licensed under the MIT License. See [LICENSE](./LICENSE) for details.