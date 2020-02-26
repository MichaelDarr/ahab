# Ahab
#### Configure, launch, and work in Dockerized environments

[![ahab](https://img.shields.io/aur/version/ahab.svg?label=ahab)](https://aur.archlinux.org/packages/ahab/) [![GitHub license](https://img.shields.io/github/license/MichaelDarr/ahab.svg)](https://github.com/MichaelDarr/ahab/blob/master/LICENSE) ![GitHub badge](https://github.com/MichaelDarr/ahab/workflows/build/badge.svg) ![GitHub badge](https://github.com/MichaelDarr/ahab/workflows/test/badge.svg) 

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

### Table of Contents
* [Installation](#installation)
* [Commands](#commands)
* [Configuration File Reference](#configuration-file-reference)
* [Key Features](#key-features)
* [FAQ](#faq)

## Installation

### Arch Linux
| AUR Package                                                   | Builds From
| :------------------------------------------------------------ | :--
| [`ahab`](https://aur.archlinux.org/packages/ahab/)            | Latest Release
| [`ahab-git`](https://aur.archlinux.org/packages/ahab-git/)    | Github Master Branch 

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
Ahab supports all Docker container commands, automatically targeting the container configured by
the `ahab.json` file located in the active directory (or in the closest parent, like `git`). For
more info on these, use the `--help` CLI option or visit the
[Docker CLI reference](https://docs.docker.com/engine/reference/commandline/cli/). Listed below are
the new commands introduced by Ahab.

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

## Configuration File Reference
**Note:** All string-type configuration fields support environment variable expansion in the forms of
`${var}` and `$var`.

### Project Configuration: `ahab.json`

| Field                 | Type          | Default       | Description
| :-------------------- | :------------ | :------------ | :--
| `ahab`                | string        | **REQUIRED**  | Minimum Ahab version required to launch this project
| `command`             | string        | `top -b`      | [See Docker Reference](https://docs.docker.com/engine/reference/run/#cmd-default-command-or-options)
| `entrypoint`          | string        | None          | [See Docker Reference](https://docs.docker.com/engine/reference/run/#entrypoint-default-command-to-execute-at-runtime)
| `environment`         | [string]      | []            | List of `KEY=VALUE` environment variables to set in the container
| `hostname`            | string        | None          | Container host name
| `image`               | string        | **REQUIRED**  | Docker image used by the container
| `init`                | [string]      | []            | List of commands to be run as root immediately after container creation
| `name`                | string        | None          | Manually assign a name to the container instead of generating it from the config path
| `options`             | [string]      | []            | List of options passed during container creation
| `permissions`         | {[permissions](#permissions)} | {}            | See [Permissions](#permissions)
| `restartAfterSetup`   | boolean       | false         | If true, container restarts after permissions are set up and `init` commands are run
| `shareX11`            | boolean       | false         | If true, container can launch windows onto the host's X11-Compatible Desktop
| `user`                | string        | `ahab`        | User for commands after initial setup
| `volumes`             | [string]      | []            | List of volumes to mount during container creation
| `workdir`             | string        | None          | [See Docker Reference](https://docs.docker.com/engine/reference/run/#workdir)

#### Permissions

| Field     | Type      | Default   | Description
| :-------- | :-------- | :-------- | :--
| `cmdSet`  | string    | None      | Commands to use while setting permissions. If using a BusyBox-based image, set to `busybox`
| `disable` | boolean   | false     | If true, Ahab does not set up any user permissions or create the `ahab` user
| `groups`  | [string]  | []        | Groups the `ahab` user belongs to. If a group does not exist in-image, prefix it with `!`

### User Configuration: `~/.config/ahab/config.json`

| Field                 | Type          | Default       | Description
| :-------------------- | :------------ | :------------ | :--
| `environment`         | [string]      | []            | List of extra `KEY=VALUE` pairs of environment variables to set in containers
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
containers on a single machine, or even in a single project.

Ahab itself has additional nested configuration files to run its test suite(`test/ahab.json`) or
build distro-specific packages(`build/deb/ahab.json`).

### Stateless Operation
Ahab has no internal state and never writes to any files. Each invocation is a blank slate, so you
never end up with a bad installation or inexplicable errors. Additionally, since Ahab is a single
compiled binary instead of an interpreted python package (like docker-compose), you don't have to
worry about machine-specific interpreter/dependency issues.

### Transparent Behavior
Ahab acts like a developer, only using Docker commands interact with containers. The commands run
by Ahab are printed directly in the shell, so you know exactly what it's doing - there's no "magic"
here. This has two major benefits:

1. It's easy to replace Ahab with other tooling if needed
2. Your own Docker CLI skills improve as you use Ahab

### Universal Image Support
With Ahab's robust configuration options, you don't need to manage custom images or Dockerfiles for
your containerized projects! You can use custom users/groups, install new packages, and mount
volumes - all without a Dockerfile or helper shell script.

## FAQ

### How is this different than docker-compose?
Ahab was created in direct response to issues with docker-compose. Docker-compose is great, but
it just isn't the right tool for the job for many projects.

One major issue with docker-compose is its lag time behind Docker itself. Docker-compose is a
python package which interfaces with Docker via docker-py, so the docker-compose team has to
wait on docker-py for feature support. At the time of writing, the `--gpus all` option is
*still* unsupported, 7 months after the feature was released. Ahab supports arbitrary launch
options, so this class of issue is entirely avoided. Additionally, it's a system-package-managed
binary instead of a python package often installed with pip, which avoids common python package
installation errors and inconsistencies.

Second, docker-compose is overkill for most projects. Its features are great if you are building
a project which needs to run multiple networked containers, but for most projects, most of what
docker-compose offers is unnecessary. Ahab is dead simple, with smaller config files and a
gentler learning curve. Its focus on transparency and replaceability make it a great option
for new projects, whose needs can change quickly.
