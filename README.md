# Bitbucket Pipelines self-hosted Runners autoscaler

This repository contains all the code related to the Bitbucket Pipelines self-hosted Runners autoscaler.

## Overview

The Bitbucket Pipelines Self-Hosted Runners Autoscaler project provides an automated solution for scaling self-hosted runners in response to pipeline workloads. By dynamically provisioning and deprovisioning runners based on job demand, this solution optimizes resource utilization and minimizes costs, while ensuring high performance and availability for Bitbucket Pipelines.

## Local Development Environment Details

### Docker

This development environment is containerized using Docker, ensuring consistency across different setups. To get started, ensure you have Docker installed on your system. If not, you can download and install Docker from [here](https://www.docker.com/products/docker-desktop/).

## Getting Started

To work on this project:

1. Ensure you have [Docker](https://www.docker.com/products/docker-desktop/) and [Visual Studio Code](https://code.visualstudio.com/download) installed.
2. Clone (or fork) this repository:
   ```shell
   git clone git@github.com:marcodellorto/bitbucket-runner-autoscaler.git
   ```
5. Open the cloned folder in Visual Studio Code.
6. VS Code should detect the DevContainer configuration. Click on the "Reopen in Container" prompt in the bottom-right corner.
7. Voila! Happy coding :)

For more detailed instructions, refer to the [DevContainer documentation](https://code.visualstudio.com/docs/devcontainers/containers).
