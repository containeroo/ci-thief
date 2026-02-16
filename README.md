# ci-thief

CLI tool to fetch GitLab CI/CD variables and print them as shell `export` lines.

## Installation

### macOS (Homebrew)

```shell
brew install containeroo/tap/ci-thief
```

## Authentication

Before using `ci-thief`, login with a GitLab personal access token:

```shell
ci-thief login --hostname gitlab.com
```

Or non-interactively:

```shell
ci-thief login --hostname gitlab.com --token "$GITLAB_TOKEN"
```

Notes:
- `--hostname` accepts values like `gitlab.com` or `https://gitlab.example.com`.
- Token scope should include `api`.
- Credentials are stored in `~/.config/ci-thief/login.json`.

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
  logout      Logout from GitLab
  version     Print the version of ci-thief

Flags:
  -h, --help            help for ci-thief
  -R, --non-recursive   Do not fetch variables from parent groups (project-only)

Use "ci-thief [command] --help" for more information about a command.
```

### Example output

```bash
ci-thief 1234
```

And you will get something like this:

```bash
# env scope parentgroup/*
export MY_OTHER_ENV_VAR='another secret value'
# env scope *
export MY_ENV_VAR='secret value'
```

### Load directly into your shell

```bash
eval "$(ci-thief 1234)"
```

By default, `ci-thief` fetches:
- Project variables from the target project.
- Group variables from the immediate group and all parent groups.

Use `--non-recursive` to fetch only project variables.
Output is ordered as parent groups -> child groups -> project, so project-level variables take precedence when evaluated in a shell.

## Security notes

- Output includes secret values in plain text. Treat output like sensitive data.
- Avoid logging command output in CI/job logs unless intentional.
