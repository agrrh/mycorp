# MyCorp

## Overview

MyCorp is an open-source Internal Developer Platform (IDP).

It provides a unified interface for developers to manage internal resources.

Operators (platform engineers) define declarative scenarios in YAML format.

Those Scenarios are then executed by the CLI tool to manage resources like virtual machines, DNS records, and databases.

### Key Features

- **Platform-first approach**: Scenarios are dynamically fetched via API and CLI, ensuring users always interact with the latest interface specifications

- **Terminal-centric design**: The CLI tool offers powerful operation capabilities through structured commands, with features like help messages, strict call formatting, and autocompletion

- **Declarative workflow**: Operators define resource management scenarios using YAML, while Users execute these scenarios through the CLI

## Main flow

```mermaid
sequenceDiagram
    box gray User expertise & visibility scope
        participant developer
        participant cli
    end

    participant developer as 👩🏻‍💻 Developer
    participant operator as 👨‍💻 Operator

    participant cli as 💻 MyCorp CLI
    participant api as ⚙️ MyCorp API

    participant iac as ☁️ IaC Provider

    %% Flow

    operator ->> api: write declarative Scenarios
    api -->> scenario: Get actual Scenario steps

    cli ->> api: requests actual scenarios
    developer ->>+ cli: "Create DNS Record"
    cli -->> api: calls Scenario run
    api -->> iac: create resource
    iac -->> api: response
    api -->> cli: returns IaC response data
    cli ->>- developer: outputs result
```

## Contribution

- See [Issues](https://github.com/agrrh/mycorp/issues) or bring your own idea
- Fork this repo
- Push to your fork's branch
- Create a Pull Request here

## Roadmap / Future plans

See [GitHub Milestones](https://github.com/agrrh/mycorp/milestones)
