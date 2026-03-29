# MyCorp

## Overview

MyCorp is an open-source Internal Developer Platform (IDP).

It provides a unified interface for developers to manage internal resources.

Operators (platform engineers) define declarative scenarios in YAML format.

Those Scenarios are then executed by the CLI tool to manage resources like virtual machines, DNS records, and databases.

```mermaid
sequenceDiagram
    box gray User expertise & visibility scope
        participant developer
        participant cli
    end

    participant developer as 👩🏻‍💻 Developer
    participant operator as 👨‍💻 Operator

    participant scenario as 📄 Scenario
    
    participant cli as 💻 CLI
    participant api as ⚙️ API

    participant iac as ☁️ IaC Provider

    %% Flow

    operator ->> scenario: write declarative Scenarios
    api -->> scenario: Get actual Scenario steps

    cli ->> api: requests actual scenarios
    developer ->>+ cli: "Create DNS Record"
    cli -->> api: calls Scenario run
    api -->> iac: create resource
    iac -->> api: response
    api -->> cli: returns IaC response data
    cli ->>- developer: outputs result
```

## Key Features

- **Platform-first approach**: Scenarios are dynamically fetched via API and CLI, ensuring users always interact with the latest interface specifications

- **Terminal-centric design**: The CLI tool offers powerful operation capabilities through structured commands, with features like help messages, strict call formatting, and autocompletion

- **Declarative workflow**: Operators define resource management scenarios using YAML, while Users execute these scenarios through the CLI

## Architecture

```mermaid
graph TD
  operator[Operator] --> |1. Define scenarios| scenarios[scenarios/]
  operator --> |2. Use CLI| cli[cmd/]
  
  user[User] --> |3. Use CLI| cli
  cli --> |4. Interact with API| api[server/]
  api --> |5. Process scenarios| controller[internal/]
  
  controller --> |6. Sync state| application[Application]
  user --> |7. Observe state| ui[cmd/]
```

## Project Structure

```
├── cmd/             # Command-line interface for Users and API server
├── internal/        # Internal models, services, and utility packages
├── pkg/             # Publicly available packages (e.g. models)
└── scenarios/       # Scenario definitions for different resource types
```

## Contribution

[GitHub Issues & Milestones](https://github.com/agrrh/mycorp/milestones)
