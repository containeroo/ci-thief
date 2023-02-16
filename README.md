# ci-thief

<!--toc:start-->

- [ci-thief](#ci-thief)
  - [Installation](#installation)
    - [macOS (Homebrew)](#macos-homebrew)
  - [Usage](#usage)
  - [Environment Variables](#environment-variables)
  <!--toc:end-->

Simple CLI tool to fetch GitLab CI variables and output them as exportable
environment variables.

## Installation

### macOS (Homebrew)

```shell
brew install containeroo/tap/ci-thief
```

## Usage

```shell
export GITLAB_HOST=git.example.com
export GITLAB_TOKEN=your-gitlab-token

ci-thief 123
```

## Environment Variables

| Variable     | Description                               |
| ------------ | ----------------------------------------- |
| GITLAB_HOST  | The hostname of your GitLab instance      |
| GITLAB_TOKEN | The GitLab API token (PAT with api scope) |
