# ahab
Dockerize your project, git style

## Objectives
Containerization is awesome, but the upfront costs of project setup and steep learning curve can
make it a pain. Ahab is a CLI tool that jump-starts this process, avoiding frustration without obfuscating your workflow.

## What does Ahab do?
Ahab searches for a project config file (`ahab.json`) and uses it to create and interact with
Docker containers. The Ahab CLI supports all of Docker's container commands without having to
manually target your container name or ID.

In addition the official Docker commands, Ahab supports new ones such as `bash`, `up`, and
`status` - we'll talk more about these and more later. Ahab provides everything you need to
quickly, easily, and effectively develop a containerized project.

## Why Ahab?

### Transparent Behaviour
Ahab acts like a developer, only using Docker commands interact with containers. The commands run
by Ahab are printed directly in the shell, so you know exactly what it's doing - there is no
"magic" here. This has two major benefits:

1. It's easy to replace Ahab with other tooling, should you need to
2. Your own Docker CLI skills improve as you use Ahab

### Container Permissions, Solved
Never worry about file permissions in your bind-mounted directories again. Ahab sets up a non-root
user inside your container with the proper credentials to create, remove, and edit bind-mounted
files while maintaining their existing permissions. It does this in a highly configurable and
transparent manner so you have complete control over your project.

### Stateless Operation
Ahab has no internal state and never writes to any files. Each invocation is a blank slate, so you
never end up with a bad installation or inexplicable errors. Additionally, since Ahab is a single
compiled binary instead of an interpreted python package (like docker-compose), you don't have to
worry about machine-specific interpreter/dependency issues.

## Installation

### Distribution Packages

#### Arch Linux
Ahab is available via the [`ahab-git`](https://aur.archlinux.org/packages/ahab-git/) AUR package

### Traditional Build

**Prerequisites**
* go (tested on 1.13.7, requires module support)

```sh
$ git clone git@github.com:MichaelDarr/ahab.git
$ cd ahab
$ make
$ make install
```

### Build Ahab *with* Ahab

**Prerequisites**
* Docker
* Ahab

```sh
$ git clone git@github.com:MichaelDarr/ahab.git
$ cd ahab
$ ahab cmd make
$ make install
```
