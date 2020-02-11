# ahab
Dockerize your project, git style

## Compile from source

### Build ahab *with* ahab
Naturally, ahab itself has an ahab configuration file! If you have Docker and `ahab` installed, you
can build ahab like so:

```sh
$ git clone git@github.com:MichaelDarr/ahab.git
$ cd ahab
$ ahab cmd make
```

### Build ahab *without* ahab
Ahab has very few dependencies, so it's easy enough to build without containerization.

**Prerequisites**
* go (>=1.11, as ahab uses go modules)

```sh
$ git clone git@github.com:MichaelDarr/ahab.git
$ cd ahab
$ make
```