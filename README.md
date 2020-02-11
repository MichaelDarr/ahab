# ahab
Dockerize your project, git style

## Compile from source

### Build ahab *with* ahab
Naturally, ahab itself can be built within an ahab-configured container!

**Prerequisites**
* Docker
* ahab

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