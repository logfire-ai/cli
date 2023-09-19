# Logfire CLI

`logfire` on the command line brings login, signup, sources, teams, stream and other logfire features to the terminal next to where you are already working with your code.

![logfire](https://github.com/logfire-sh/cli-private/assets/28614457/ff057447-c898-47a0-ae32-529066ce57db)

## Features

You can run the following commands from the command-line interface (CLI) by directly passing arguments. If your terminal is interactive, you have the option to skip providing the arguments upfront and provide them during runtime.

- login
- signup
- logout
- teams (list, create, list-members)
- sources (list, create, delete)
- stream (livetail)

## Installation

### macOS and linux

`logfire` is available on **macOS** and **linux** via [Homebrew](https://brew.sh/), [APT](https://ubuntu.com/server/docs/package-management), [YUM](http://yum.baseurl.org/).

```bash
$ brew tap logfire-sh/tap
$ brew install logfire
```

```bash
$ curl -O https://raw.githubusercontent.com/logfire-sh/cli/main/add_logfire_cli_yum.sh
$ chmod +x add_logfire_cli_yum.sh
$ ./add_logfire_cli_yum.sh
$ yum install logfire
```

```bash
$ curl -s --compressed "https://logfire-sh.github.io/cli/KEY.gpg" | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/KEY.gpg >/dev/null
$ sudo curl -s --compressed -o /etc/apt/sources.list.d/my_list_file.list "https://logfire-sh.github.io/cli/my_list_file.list"
$ sudo apt-get update
$ sudo apt install logfire
```

### Windows

`logfire` is available on **windows** available via [Scoop](https://scoop.sh/).

```bash
$ scoop bucket add logfire-ai https://github.com/logfire-sh/cli.git
$ scoop install logfire-ai/logfire
```

## Release

We are using [goreleaser](https://goreleaser.com/) to automate creation of release artifacts for various operating systems and architectures.

GitHub actions will automatically handle syncing and publishing.
