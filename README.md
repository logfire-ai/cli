# Logfire CLI

`logfire` on the command line brings login, signup, sources, teams, stream and other logfire features to the terminal next to where you are already working with your code.

![logfire](https://github.com/logfire-sh/cli-private/assets/28614457/ff057447-c898-47a0-ae32-529066ce57db)

## Features

You can run the following commands from the command-line interface (CLI) by directly passing arguments. If your terminal is interactive, you have the option to skip providing the arguments upfront and provide them during runtime.

- login
- signup
- teams (list)
- sources (list, create, delete)
- stream (livetail)

## Installation

### macOS and linux

`logfire` is available on **macOS** and **linux** via [Homebrew][].

```bash
$ brew tap logfire-sh/tap
$ brew install logfire
```

### Windows

`logfire` is available on **windows** available via [scoop][].

```bash
$ scoop bucket add logfire-sh https://github.com/logfire-sh/cli.git
$ scoop install logfire-sh/logfire
```

## Release

We are using [goreleaser](https://goreleaser.com/) to automate creation of release artifacts for various operating systems and architectures.

To do a release run the following commands after pushing the code to master branch:

```bash
$ git tag -a <tag-name> -m <tag-message>
$ git push origin <tag-name>
```

CI/CD will autoamtically handle syncing and publishing the artifacts.
