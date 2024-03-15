![The devzero logo](https://console.devzero.io/_next/image?url=%2F_next%2Fstatic%2Fmedia%2Ffull_logo.379052d8.png&w=1080&q=75)

# LDA

The **LDA** CLI is local developer analytics tool. It provides:

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

Install from source:

- Clone the project
- Change directory
- Run `make install`
  - binary will be in your `$GOPATH/bin`
  - if your PATH isn't set correctly: export GOPATH=$(go env GOPATH) && export PATH=$PATH:$GOPATH/bin

## Install & Usage

Lda's help interface provides summaries for commands and flags:

```sh
lda --help
```

After binary has been compiled we can install the service with the following commands:

* `lda install` => This will install the daemon, configure base directory, and inject configuration into the shell
* `lda start` => This will start the daemon
* `lda stop` => This will stop the daemon
* `lda uninstall` => This will uninstall the LDA and remove all configuration
* `lda serve` => This will serve the local dashbaord with data overview

## Community

For help and discussion around LDA, best practices, and more, join us on [Slack][badges_slack].

For updates on the LDA CLI, [follow this repo on GitHub][repo].

For feature requests, bugs, or technical questions, email us at <>.
