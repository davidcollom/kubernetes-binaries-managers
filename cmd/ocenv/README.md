![Static Tests](https://github.com/little-angry-clouds/kubernetes-binaries-managers/workflows/Generic%20tests/badge.svg) ![Int Test Linux](https://github.com/little-angry-clouds/kubernetes-binaries-managers/workflows/Int%20Test%20Linux/badge.svg) ![Int Test MacOS](https://github.com/little-angry-clouds/kubernetes-binaries-managers/workflows/Int%20Test%20MacOS/badge.svg) ![Int Test Windows](https://github.com/little-angry-clouds/kubernetes-binaries-managers/workflows/Int%20Test%20Windows/badge.svg)

# ocenv

[OC](https://docs.openshift.com/container-platform/4.2/cli_reference/openshift_cli/getting-started-cli.html) version
manager inspired by [tfenv](https://github.com/tfutils/tfenv/).

## Features

- Install OC versions in a reproducible and easy way
- Enforce version in your git repositories with a `.oc_version` file

## Supported OS

Currently ocenv supports the following OSes

- Mac OS
- Linux
- Windows

## Installation

There are two components in `ocenv`. One is the `ocenv` binary, the other one
is a `oc` wrapper. It works as if were `oc`, but it has some logic to choose
the version to execute. You should take care and ensure that you don't have any
`oc` binary in your path. To check which binary you're executing, you can see
it with:

```bash
$ which oc
/opt/brew/bin/oc
```

### Homebrew

This is the recomended way, since it provides upgrades. It should work in Mac,
Linux and Windows with WSL.

```bash
# Just the first time, activate the repository
brew tap little-angry-clouds/homebrew-my-brews
# To install
brew install ocenv
# To upgrade
brew upgrade ocenv
```

You should add your `homebrew` binary path to your PATH:

```bash
echo 'export PATH="$(brew --prefix)/bin/:$PATH"' >> ~/.bashrc
# Or
echo 'export PATH="$(brew --prefix)/bin/:$PATH"' >> ~/.zshrc
```

For Windows you should do the weird stuff that it needs to set an environmental variable.

### Manually

1. Add `~/.bin` to your `$PATH` and create it if doesn't exist

```bash
echo 'export PATH="$HOME/.bin:$PATH"' >> ~/.bashrc
# Or
echo 'export PATH="$HOME/.bin:$PATH"' >> ~/.zshrc

mkdir -p ~/.bin
```

2. Download the binaries and put them on your path

Go to [the releases
page](https://github.com/little-angry-clouds/kubernetes-binaries-managers/releases)
and download the version you want. For example:

```bash
wget https://github.com/little-angry-clouds/kubernetes-binaries-managers/releases/download/v1.0.0/kubernetes-binaries-managers_1.0.0_linux_amd64.tar.gz
tar -xzf kubernetes-binaries-managers_1.0.0_linux_amd64.tar.gz
mv oc-linux-amd64/ocenv ~/.bin/ocenv
mv oc-linux-amd64/oc-wrapper ~/.bin/oc
```

And that's it!

## Usage

### Help

```bash
$ ocenv help
Oc version manager

Usage:
  ocenv [command]

Available Commands:
  help        Help about any command
  install     Install binary
  list        Lists local and remote versions
  uninstall   Uninstall binary
  use         Set the default version to use

Flags:
  -h, --help     help

Use "ocenv [command] --help" for more information about a command.
```

### List installable versions

This option uses Github API to paginate all versions. Github API has some usage
limitations. It usually works, but if you happen to do a lot of requests to
github or are on an office or similar, chances are that this command will fail.
You can still install binaries if you know the version you want, thought.

```bash
$ ocenv list remote
4.15.0-0.okd-2024-03-10-010116
4.15.0-0.okd-2024-02-23-163410
4.15.0-0.okd-2024-02-10-035534
4.15.0-0.okd-2024-01-27-070424
4.14.0-0.okd-2024-01-26-175629
4.14.0-0.okd-2024-01-06-084517
4.14.0-0.okd-2023-12-01-225814
4.14.0-0.okd-2023-11-14-101924
4.14.0-0.okd-2023-11-12-042703
4.14.0-0.okd-2023-10-28-073550
...
```

### List installed versions

```bash
$ ocenv list local
4.15.0-0.okd-2024-03-10-010116
4.15.0-0.okd-2024-02-10-035534
```

### Install version

```bash
$ ocenv install 4.14.0-0.okd-2024-01-06-084517
Downloading binary...
Done! Saving it at /home/user/.bin/oc-4.14.0-0.okd-2024-01-06-084517
```

### Use version

```bash
$ ocenv use 4.14.0-0.okd-2024-01-06-084517
Done! Using 4.14.0-0.okd-2024-01-06-084517 version.
```

### Uninstall version

```bash
$ ocenv uninstall 4.14.0-0.okd-2024-01-06-084517
Done! 4.14.0-0.okd-2024-01-06-084517 version uninstalled from /home/ap/.bin/oc-4.14.0-0.okd-2024-01-06-084517.
```

## License

GPL3
