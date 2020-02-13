#/usr/bin/env bash
# This script creates a .deb package using the ahab executable in the repo's root
# Run this from within an ubuntu container using ahab (assuming cwd=this dir)
#   $ ahab bash
#   $ ./pack.sh
set -e

# Ensure this script is being run from this dir, not caller's
cd "$(dirname "$0")"

SRCDIR=$PWD/../..

# Set up debian-style package name
EXECUTABLE=$SRCDIR/ahab
REVISION=1
VERSION="$(cat $SRCDIR/VERSION)"
PKGVERSION=$VERSION-$REVISION
PKGNAME=ahab_$PKGVERSION

# Ensure ahab executable exists
if [ ! -f "$EXECUTABLE" ]; then
    echo "Could not find Ahab executable - is it built?"
    exit 1
fi

# Create and populate package directory structure
mkdir -p $PKGNAME/DEBIAN $PKGNAME/usr/local/bin
cp $EXECUTABLE $PKGNAME/usr/local/bin/ahab
echo "\
Package: ahab
Version: $PKGVERSION
Section: base
Priority: optional
Architecture: amd64
Recommends: docker
Maintainer: Michael Darr <michael.e.darr@gmail.com>
Description: Dockerize your project, git style
Homepage: https://github.com/MichaelDarr/ahab" > $PKGNAME/DEBIAN/control

# Build the package
dpkg-deb --build $PKGNAME

# Remove build artifacts, leaving only the .deb file
rm -r $PKGNAME
