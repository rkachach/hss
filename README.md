# HTTP Storage Service (hss)

[![License](https://img.shields.io/badge/license-MPL--2.0-brightgreen.svg)](https://github.com/rkachach/hss/blob/main/LICENSE)
[![GitHub Issues](https://img.shields.io/github/issues/your-username/http-storage-service.svg)](https://github.com/rkachach/hss/issues)
[![GitHub Stars](https://img.shields.io/github/stars/your-username/http-storage-service.svg)](https://github.com/rkachach/hss/stargazers)

# Filesystem Service with REST API

## Overview

Welcome to the HTTP Storage Service!

This project provides a filesystem-like service that allows users to interact with files and directories using a RESTful API over HTTP. The service offers a convenient and platform-independent way to
manage files and directories remotely, making it suitable for various applications and environments.

## Features

- **Standardization**: Adheres to well-defined RESTful principles, promoting interoperability and ease of use.
- **Integration**: Easily integrates with other web services, applications, and automation workflows.
- **Remote Access**: Access and manipulate files and directories from anywhere with an internet connection.
- **Scalability**: Can be scaled horizontally to handle increased workload and traffic.
- **Security**: Supports HTTPS for secure communication and provides authentication and authorization mechanisms.
- **Flexibility**: Allows for custom endpoints and operations tailored to specific use cases.

## Advantages

- **Platform Agsnostic**: Easy interaction with the service regardless of platform or programming language.
- **Cross-Platform Accessibility:** Access and manage your data from any device with internet connectivity.
- **RESTful Interface:** Intuitive and user-friendly interface for seamless integration and utilization.
- **Hierarchical Organization:** Organize your data into a structured hierarchy for easier management.
- **Security Measures:** Robust security protocols ensure the safety and integrity of your stored data.
- **Developer-Friendly:** Comprehensive documentation, SDKs, and support for quick integration.
- **Remote Collaboration**: Facilitates collaboration by enabling users to share and work on files remotely.
- **Automation**: Supports automation of file-related tasks through programmatic interaction with the API.

## Getting Started

## Installation

To get started with the filesystem service, follow these steps:

1. Clone the repository: `git clone https://github.com/rkachach/hss.git`
2. Install dependencies: `go mod tidy`
3. Modify `config/config.json` specifying the directory to serve.
4. Start the service: `go run cmd/app/main.go`
5. Use test client `clients/web-client/index.html`

## Usage

For detailed API documentation, refer to the [API Reference](api/openapi.yaml).

## API Reference

For detailed documentation on the HTTP Filesystem Service API, including endpoint descriptions and examples, see the [API Reference](api/openapi.yaml).

## No Warranty

This software is provided as-is without any warranty. Use it at your own risk. The authors and contributors of this software are not liable for any damages or losses arising from its use.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvement, please open an issue or submit a pull request.
