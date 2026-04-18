# AGENTS.md - Mycorp Project Context

## Project Overview

**Mycorp** is an open-source IDP (Internal Developer Platform).

It aims to provide unified interface for Users (developers) to operate common Scenarios.

For example, requesting or managing internal resource, such as virtual machines, DNS records, databases, etc.

Those Scenarios are owned by Operators (platform engineers) and managed as declarative specs, written in YAML.

### Key Features

- **Platform approach**: Scenarios provided via API, CLI tool always fetch latest and most stable interface spec.
- **Terminal-first**: CLI tool means simple yet powerful way to operate on tasks: help messages, strict calls format, autocompletion.

## Project Structure

```
├── cmd
│   ├── client                # Command-line interface for Users
│   └── server                # REST API which serve and run scenarios written by Operators
├── internal
│   ├── application
│   │   └── server            # REST API Handlers
│   ├── domain                # Domain business logic
│   │   └── modules           # Specific modules
│   └── infrastructure        # Database connectors etc.
├── scenarios                 # Scenario examples
└── tmp
```

## Running

### Server / API

```sh
SCENARIO_DIR="$(pwd)/scenarios" air --build.cmd "go build -o tmp/server cmd/server/main.go" --build.entrypoint "./tmp/server"
```

### CLI

Expected API is available at [http://127.0.0.1:8080](http://127.0.0.1:8080):

```sh
MYCORP_CLI_SCENARIOS_URL="http://127.0.0.1:8080/scenarios" go run ./cmd/client dns create --zone example.org
```
