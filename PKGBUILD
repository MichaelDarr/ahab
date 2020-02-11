# Maintainer: Michael Darr <michael.e.darr@gmail.com>
_pkgname="ahab"
pkgname="${_pkgname}-git"
pkgver=0.1
pkgrel=1
pkgdesc="Dockerize your project, git style"
arch=('x86_64' 'aarch64')
url="https://github.com/MichaelDarr/ahab"
license=('GPL')
depends=(
  'git'
  'sudo')
makedepends=('go')
conflicts=('ahab')
provides=('ahab')
source=("ahab::git+https://github.com/MichaelDarr/ahab.git#branch=master")
md5sums=("SKIP")

pkgver() {
	cat "$srcdir/$_pkgname/VERSION"
}

build() {
	cd "$srcdir/$_pkgname"
	make build
}

package() {
	cd "$srcdir/$_pkgname"
	make DESTDIR="$pkgdir" PREFIX=/usr install
}
