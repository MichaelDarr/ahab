# Ahab
#### Configure, launch, and work in Dockerized environments
![GitHub badge](https://github.com/MichaelDarr/ahab/workflows/build/badge.svg) ![GitHub badge](https://github.com/MichaelDarr/ahab/workflows/test/badge.svg)

[![Go Report Card](https://goreportcard.com/badge/github.com/MichaelDarr/ahab)](https://goreportcard.com/badge/github.com/MichaelDarr/ahab) [![Godoc Reference](https://godoc.org/github.com/MichaelDarr/ahab/internal?status.svg)](https://godoc.org/github.com/MichaelDarr/ahab/internal) [![GitHub license](https://img.shields.io/github/license/MichaelDarr/ahab.svg)](https://github.com/MichaelDarr/ahab/blob/master/LICENSE)

Containerization is awesome, but the upfront costs of project setup and steep learning curve can
make it a pain. Ahab is a CLI tool that jump-starts this process, avoiding frustration without
obfuscating your workflow.

When invoked at the command line, `ahab` searches the current directory (and its parents,
recursively) for a file `ahab.json`. Its contents describe the manner in which Ahab will create and
interact with a project's Docker container. The Ahab CLI supports all Docker container commands
without having to specify a container name or ID - the container described by `ahab.json` is
automatically targeted. For example, `docker rm [container_id]` becomes `ahab rm`. Additionally,
Ahab introduces new commands such as `bash`, `up`, and `status` - a full list can be found
[here](#commands). Ahab provides everything you need to quickly and effectively develop a
containerized project.

## Table of Contents
* [Installation](#installation)
* [Commands](#commands)
* [Configuration Reference](#configuration-reference)
* [Key Features](#key-features)
* [FAQ](#faq)

## Installation

### Arch Linux
| AUR Package | Builds From
| :-- | :--
| [![ahab](https://img.shields.io/aur/version/ahab.svg?label=ahab)](https://aur.archlinux.org/packages/ahab/) | Latest Release
| [![ahab-git](https://img.shields.io/aur/version/ahab-git.svg?label=ahab-git)](https://aur.archlinux.org/packages/ahab-git/)    | Github Master Branch 

### Ubuntu (bionic, disco, eoan, xenial)
```sh
$ sudo add-apt-repository ppa:michaeldarr/ppa
$ sudo apt update
$ sudo apt install ahab
```

### Other (From Source)

**Prerequisites:**
* git
* go (tested on 1.13.7)
* make

```sh
$ git clone git@github.com:MichaelDarr/ahab.git
$ cd ahab
$ make  # or, `make self` to build ahab *with* ahab (prerequisites: Docker, Ahab)
$ make install
```

## Commands
Ahab supports all Docker container commands as described [here](#Ahab). For more information, use
`ahab [command] --help` or visit the
[Docker CLI reference](https://docs.docker.com/engine/reference/commandline/cli/). Listed below are
supplemental commands introduced by Ahab.

| Command               | Description
| :-------------------- | :--
| `bash`, `sh`, `zsh`   | Shell-specific container terminal access
| `cmd`                 | Execute an in-container command, attaching the input/output to your terminal
| `down`                | Stop and remove the container
| `ls`                  | List all containers, images, and volumes on your machine
| `ls{c,i,v}`           | List one Docker asset type, e.g. `lsc`ontainers
| `prune`               | Remove all unused Docker assets on your machine
| `status`              | Print a human-friendly container status report
| `up`                  | Create and start the container

## Configuration Reference
**Note:** All string-type configuration fields support environment variable expansion in the forms
`${var}` and `$var`.

### Project Configuration: `ahab.json`

| Field                 | Type          | Default       | Description
| :-------------------- | :------------ | :------------ | :--
| `ahab`                | string        | **REQUIRED**  | Minimum Ahab version required to launch a project
| `buildContext`        | string        | `ahab.json`'s directory         | Docker build context (if using `dockerfile` instead of `image`)
| `command`             | string        | `top -b`      | [See Docker Reference](https://docs.docker.com/engine/reference/run/#cmd-default-command-or-options)
| `dockerfile`          | string        | None          | Dockerfile used to build container image (required if `image` not present)
| `entrypoint`          | string        | None          | [See Docker Reference](https://docs.docker.com/engine/reference/run/#entrypoint-default-command-to-execute-at-runtime)
| `environment`         | [string]      | []            | List of `KEY=VALUE` environment variables to set in the container
| `hostname`            | string        | None          | Container host name
| `image`               | string        | None          | Docker image used by the container  (required if `dockerfile` not present)
| `init`                | [string]      | []            | List of commands to be run as root immediately after container creation
| `name`                | string        | None          | Manually assign a name to the container instead of generating it from the config path
| `options`             | [string]      | []            | List of options passed during container creation
| `permissions`         | {[permissions](#permissions)} | {}            | See [Permissions](#permissions)
| `restartAfterSetup`   | boolean       | false         | If true, the container restarts after permissions are set up and `init` commands are run
| `shareX11`            | boolean       | false         | If true, processes within the container can launch windows onto the host's X11-compatible desktop
| `user`                | string        | `ahab`        | User for container commands after initial setup
| `volumes`             | [string]      | []            | List of volumes to mount during container creation
| `workdir`             | string        | None          | [See Docker Reference](https://docs.docker.com/engine/reference/run/#workdir)

#### Permissions

| Field     | Type      | Default   | Description
| :-------- | :-------- | :-------- | :--
| `cmdSet`  | string    | None      | Commands to use while setting permissions. If using a BusyBox-based image, set to `busybox`
| `disable` | boolean   | false     | If true, Ahab does not set up any user permissions or create the `ahab` user
| `groups`  | [string]  | []        | Groups the `ahab` user belongs to. If a group does not yet exist in-image, prefix it with a `!`, e.g. `"groups": ["!docker"]`

### User Configuration: `~/.config/ahab/config.json`

| Field                 | Type          | Default       | Description
| :-------------------- | :------------ | :------------ | :--
| `environment`         | [string]      | []            | List of extra `KEY=VALUE` environment variables to set in containers
| `hideCommands`        | boolean       | false         | If true, Ahab will not print Docker commands before it runs them
| `options`             | [string]      | []            | List of extra options passed during container creation
| `volumes`             | [string]      | []            | List of extra of volumes to mount during container creation

## Key Features

### Container Permissions, Solved
Never worry about file permissions in your bind-mounted directories again. Ahab sets up a non-root
user inside your container with the proper credentials to create, remove, and edit bind-mounted
files while maintaining their existing permissions. It does this in a highly configurable and
transparent manner so you have complete control over your project.

### Git-Style Configuration
Ahab looks for the `ahab.json` file like `git` looks for the `.git` folder. If there is no config
file in the current directory, Ahab searches its parent directories recursively until `ahab.json`
is found or the root of the file system is reached. With this system, it's easy to manage multiple
containers on a single machine, or even in a single project. Internally, Ahab leverages this
feature to run its test suite (`test/ahab.json`) and build distribution-specific packages
(`build/deb/ahab.json`).

### Stateless Operation
Ahab has no internal state and never writes to any files. Each invocation is a blank slate, so you
never end up with a bad installation or inexplicable errors. Additionally, since Ahab is an
independent compiled binary instead of an interpreted python package (like `docker-compose`), you
don't have to worry about machine-specific interpreter/dependency issues.

### Transparent Behavior
Ahab acts like a developer, only using Docker commands interact with containers. The commands run
by Ahab are printed directly in the shell, so you know exactly what it's doing - there's no "magic"
here. This has two major benefits:

1. It's easy to replace Ahab with other tooling if needed
2. Your own Docker CLI skills improve as you use Ahab

### Universal Image Support
With Ahab's robust configuration options, you don't need to manage custom images or Dockerfiles for
your containerized projects! You can set up users/groups, install new packages, and utilize arbitrary
container options all without a single Dockerfile or shell script.

## FAQ

### How is this different than docker-compose?
Ahab was created in response to frustration with docker-compose. One issue with docker-compose is
its lag time behind Docker itself. Since docker-compose is a python package which interfaces with
Docker via docker-py, the docker-compose team has to wait on docker-py for feature support. At the
time of writing, the `--gpus` option is unsupported by docker-compose 7 months after its release in
Docker CE. Ahab's arbitrary container options entirely avoid this issue.

A second concern over docker-compose is its distribution mechanism as a python package. For teams
of developers, this can be problematic due to inconsistent python package management and dependency
tree conflicts. Internal tooling to extend docker-compose further exacerbates these issues. In
contrast, Ahab is installed as a standalone binary executable.

Lastly, docker-compose and Ahab serve two very different project archetypes. The ability to run
multiple interdependent containers is needed only by a minority of projects; for most, only a
single well-configured environment is required. Here, Ahab thrives with its simpler config files,
development-ready environments, and gentler learning curves. New, volatile projects especially
benefit from Ahab's commitment to transparency and replaceability.
