# qseal

<img src="resources/gopher.png" alt="Gopher" width="200"/>


**qseal** is a CLI tool that simplifies the process of sealing and unsealing Kubernetes secrets using [`kubeseal`](https://github.com/bitnami-labs/sealed-secrets). It uses a declarative configuration file (`qsealrc.yaml`) to manage your secrets.

## Features

- Declarative configuration of sealed secrets
- Sealing and unsealing via a single `sync` operation
- Conflict detection for sealed paths

## Installation

Clone the repository and build the binary:

```bash
go install gitlab.42paris.fr/froz/qseal@latest
```

You can also download the latest release from the [releases page](https://gitlab.42paris.fr/froz/qseal/-/releases).

## Usage

```bash
qseal [flags]
qseal [command]
qseal # without any command will run `qseal sync`
```

### Available Commands

| Command      | Description                                                                                    |
| ------------ | ---------------------------------------------------------------------------------------------- |
| `init`       | Initialize the `qsealrc.yaml` configuration file                                               |
| `seal-all`   | Seal all secrets defined in the config file _(not recommended, use `qseal sync` or `qseal`)_   |
| `unseal-all` | Unseal all secrets defined in the config file _(not recommended, use `qseal sync` or `qseal`)_ |
| `completion` | Generate autocompletion script for your shell                                                  |
| `help`       | Display help for any command                                                                   |

### Flags

- `-h`, `--help`: Show help information

Use `qseal [command] --help` for detailed information about a specific command.

## Configuration

qseal expects a `qsealrc.yaml` file at the root of your project. This file defines all secrets to be managed. Each secret must include:

- A name
- The path to the sealed file
- Then the path to the secret file (env file, files)
- The type of secret (e.g., `kubernetes.io/dockerconfigjson`, `kubernetes.io/tls`, etc.)

## Sync Logic

The core of `qseal` is the `Sync` operation, which:

1. Parses the secrets listed in `qsealrc.yaml`
2. Groups them by sealed output path
3. Determines whether each group needs to be sealed, unsealed, or skipped
4. Detects conflicts (e.g. multiple actions for the same sealed path)
5. Applies sealing or unsealing as needed

Example log output:

```txt
[2025-04-16 10:00:00] SEALING secrets.yaml (3 secret(s))
[2025-04-16 10:00:00] SKIP secrets.yaml (up-to-date)
[2025-04-16 10:00:00] UNSEALING secrets.yaml (2 secret(s))
```

## Conflict Handling

If multiple secrets reference the same sealed file path but require different actions (`seal` vs `unseal`), `qseal` will raise an error. You can resolve this by explicitly running either:

```bash
qseal seal-all
```

or

```bash
qseal unseal-all
```
