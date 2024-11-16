# Golang DevContainer

This repository serves as a skeleton for setting up Golang development environments using Visual Studio Code's DevContainers feature.

## Overview

The Golang DevContainer provides a pre-configured environment for Golang development, leveraging Visual Studio Code's remote container capabilities. It includes essential configurations and extensions to streamline the Golang development workflow.

## Features

- **Golang Environment**: Configured with the latest stable version of Golang.
- **Visual Studio Code Integration**: Optimized settings and extensions for seamless integration with VS Code.
- **Docker Configuration**: Ready-to-use Dockerfile and docker-compose configuration for easy containerization.
- **AWS and SAM Integration**: AWS Command Line Interface (CLI) and AWS Serverless Application Model (SAM) are installed to facilitate AWS-related development tasks.

## Development Environment Details

### Docker

This development environment is containerized using Docker, ensuring consistency across different setups. To get started, ensure you have Docker installed on your system. If not, you can download and install Docker from [here](https://www.docker.com/products/docker-desktop/).

### AWS CLI and SAM

The AWS CLI and SAM tools are pre-installed in this environment to facilitate AWS-related tasks. Here's how to access them within the DevContainer:

1. **AWS CLI**: The AWS CLI is configured and ready to use. If you did not yet set up your AWS credentials, run `aws configure` and start using the CLI for managing AWS services.

2. **AWS SAM**: AWS SAM is installed to simplify the development of serverless applications. Use `sam --version` to confirm the installation and explore available commands for building and deploying serverless applications.

For more information on using Docker, AWS CLI, and SAM, refer to their respective documentation:

- [Docker Documentation](https://docs.docker.com/)
- [AWS CLI Documentation](https://docs.aws.amazon.com/cli/)
- [AWS SAM Documentation](https://docs.aws.amazon.com/serverless-application-model/)

## Getting Started

To use this DevContainer for your Golang projects:

1. Ensure you have [Docker](https://www.docker.com/products/docker-desktop/) and [Visual Studio Code](https://code.visualstudio.com/download) installed.
2. Clone (or fork) this repository:
   ```shell
   git clone git@github.com:marcodellorto/golang-devcontainer.git
   ```
3. Create your `docker-compose.override.yml` file starting from the `docker-compose.override.yml.dist` file already present in the repository
   ```shell
   cp docker-compose.override.yml.dist docker-compose.override.yml
   ```
4. Open the cloned folder in Visual Studio Code.
5. VS Code should detect the DevContainer configuration. Click on the "Reopen in Container" prompt in the bottom-right corner.
6. Voila! You're now set up with a Golang development environment with Docker, AWS CLI, and SAM configured.

For more detailed instructions, refer to the [DevContainer documentation](https://code.visualstudio.com/docs/devcontainers/containers).

## Contributing

Contributions are welcome! Feel free to open issues or pull requests for any improvements or new features.

Before contributing, please review the [contribution guidelines](CONTRIBUTING.md) to maintain consistency and quality.

## License

This project is licensed under the [MIT License](LICENSE.md).

## Acknowledgments

- **Contributors**:
  - [@marcodellorto](https://github.com/marcodellorto)

- **Libraries and Tools**:
  - [Docker](https://www.docker.com/): Used for containerization of the development environment.
  - [AWS CLI](https://aws.amazon.com/cli/): Facilitates interaction with AWS services.
  - [AWS SAM](https://aws.amazon.com/serverless/sam/): Aiding in the development of serverless applications.
  - [VS Code DevContainers](https://code.visualstudio.com/docs/remote/containers): Structure and setup inspiration.
  - [golangci-lint](https://golangci-lint.run/): Linting tool for checking Go code.

- **Docker Images**:
  - [marcodellorto/golang](https://hub.docker.com/r/marcodellorto/golang): Docker images used as a base for the Golang development environment.

- **External Resources**:
  - [Golang Documentation](https://golang.org/doc/): Best practices and language features.
  - [Golang Project Layout](https://github.com/golang-standards/project-layout): Basic layout for Go application projects, not an official standard but highly suggested.
