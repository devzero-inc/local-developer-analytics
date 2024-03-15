![The devzero logo](https://assets-global.website-files.com/659f77ad8e06050cc27ed4d3/65aaf7abd1a9b99456f6154a_Devzero%20logo%20on%20dark-p-500.png)

# LDA

[![License]()]()
[![Release]()]()
[![CI]()]()
[![Homebrew]()]()

The [`lda`][lda] CLI is local developer analytics tool. It provides:

- Tracking commands executed in user shell
- Tracking processes running on users computer
- Overview of command execution
- Overview of most intesive processes over time

## Installation

### Homebrew

You can install `lda` using [Homebrew][brew] (macOS or Linux):

```sh
brew install <>
```

### Other methods

Install from binary:

- Clone the project
- Change directory
- Run `make build`
- Move binary to desired `/bin/` location

## Install & Usage

Lda's help interface provides summaries for commands and flags:

```sh
lda --help
```

After binary has been compiled we can install the service with the following commands:

`lda install` => This will install the daemon, configure base directory, and inject configuration into the shell
`lda start` => This will start the daemon
`lda stop` => This will stop the daemon
`lda uninstall` => This will uninstall the LDA and remove all configuration
`lda serve` => This will serve the local dashbaord with data overview

## Community

For help and discussion around LDA, best practices, and more, join us on [Slack][badges_slack].

For updates on the LDA CLI, [follow this repo on GitHub][repo].

For feature requests, bugs, or technical questions, email us at <>.
