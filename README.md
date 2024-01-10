# ci-thief

Simple CLI tool to fetch GitLab CI variables and output them as exportable
environment variables.

## Installation

### macOS (Homebrew)

```shell
brew install containeroo/tap/ci-thief
```

## Usage

```text
Fetch GitLab CI variables from a project

Usage:
  ci-thief [PROJECT_ID] [flags]
  ci-thief [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  login       Login to GitLab
  version     Print the version of ci-thief

Flags:
  -h, --help            help for ci-thief
  -R, --non-recursive   Do not fetch variables from parent groups

Use "ci-thief [command] --help" for more information about a command.
```
