# Ahab
Containerization is awesome, but the upfront costs of project setup and steep learning curve can
make it a pain. Ahab is a CLI tool that jump-starts this process, avoiding frustration without obfuscating your workflow.

Ahab searches for a project config file `ahab.json` and uses it to create and interact with
Docker containers. The Ahab CLI supports all of Docker's container commands without having to
manually target your container name or ID.

In addition the official Docker commands, Ahab supports new ones such as `bash`, `up`, and
`status` - we'll talk about these and more later. Ahab provides everything you need to
quickly and effectively develop a containerized project.

### Table of Contents
* [Install](#install)
* [Commands](#commands)
* [Key Features](#key-features)
* [Build From Source](#build-from-source)

## Install

### Arch Linux
| AUR Package                                                   | Builds From
| :------------------------------------------------------------ | :--
| [`ahab`](https://aur.archlinux.org/packages/ahab/)            | Latest Release
| [`ahab-git`](https://aur.archlinux.org/packages/ahab-git/)    | Github Master Branch 

### Debian-Based (dpkg)
```sh
$ wget https://github.com/MichaelDarr/ahab/releases/download/0.1/ahab_0.1-1.deb
$ sudo dpkg -i ahab_0.1-1.deb
```

### Other
Try [building from source](#build-from-source)!

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

## Config File Reference

### Project Configuration
All project configs are named `ahab.json`

| Field                 | Type          | Default       | Description
| :-------------------- | :------------ | :------------ | :--
| `ahab`                | string        | **REQUIRED**  | Minimum Ahab version required to launch this project
| `command`             | string        | `top -b`      | [See Docker Reference](https://docs.docker.com/engine/reference/run/#cmd-default-command-or-options)
| `entrypoint`          | string        | None          | [See Docker Reference](https://docs.docker.com/engine/reference/run/#entrypoint-default-command-to-execute-at-runtime)
| `environment`         | [string]      | []            | List of `KEY=VALUE` pairs of environment variables to set in the container
| `hostname`            | string        | None          | Container host name
| `image`               | string        | **REQUIRED**  | Docker image used by the container
| `init`                | [string]      | []            | List of commands to be run as root immediately after container creation
| `name`                | string        | None          | Manually assign a name to the container instead of generating it from the config path
| `options`             | [string]      | []            | List of options passed during container creation
| `permissions`         | PERMISSIONS   | None          | Permissions options - see the [Permissions](#permissions) section
| `restartAfterSetup`   | boolean       | false         | If true, container restarts after permissions are set up and `init` commands are run
| `shareX11`            | boolean       | false         | If true, container can launch windows onto the host's X11-Compatible Desktop
| `user`                | string        | `ahab`        | User used to run commands run by ahab after initial setup by root
| `volumes`             | [string]      | []            | List of volumes to mount during container creation
| `workdir`             | string        | None          | [See Docker Reference](https://docs.docker.com/engine/reference/run/#workdir)

#### Permissions

| Field     | Type      | Default   | Description
| :-------- | :-------- | :-------- | :--
| `cmdSet`  | string    | None      | Commands to use while setting permissions. If using a BusyBox-based image, set to `busybox`
| `disable` | boolean   | false     | If true, Ahab does not set up any user permissions or create the `ahab` user
| `groups`  | [string]  | []        | Groups the `ahab` user belongs to. If a group does not exist in-image, prefix it with `!`

### User Configuration
All user configs are should be located at `~/.config/ahab/config.json`

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

## Build from Source

**Prerequisites:**
* go (tested on 1.13.7)

```sh
$ git clone git@github.com:MichaelDarr/ahab.git
$ cd ahab
$ make
$ make install
```

### Build Ahab *with* Ahab

**Prerequisites:**
* Docker
* Ahab

```sh
$ git clone git@github.com:MichaelDarr/ahab.git
$ cd ahab
$ ahab cmd make
$ make install
```
