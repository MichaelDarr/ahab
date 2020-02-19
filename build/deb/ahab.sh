#!/usr/bin/env bash
# Package ahab for debian
# Usage: ./ahab.sh [VERSION]
set -e

# Ensure this script is being run from this dir, not caller's
cd "$(dirname "$0")"

# Ensure version arg is passed
if [ -z "$1" ]; then
    echo "Package ahab for Debian.
Usage: ./ahab.sh [VERSION]"
    exit 1
fi

# download & extract tarball
TARBALL=v"$1".tar.gz
EXTRACTED=ahab-"$1"
ARCHIVE=https://github.com/MichaelDarr/ahab/archive/"$TARBALL"
wget "$ARCHIVE"
tar -xzf "$TARBALL"

cd "$EXTRACTED"
dh_make -f ../"$TARBALL"

echo "Additional Steps:
* Remove template (.ex) files
* debian/changelog
  - unstable -> eoan
  - If needed, change revision number: ahab (VERSION-REVISION)
* See files in template/* for more changes

Build Package:
$ dpkg-buildpackage -S

Sign Package:
$ debsign -k [KEY ID] [PACKAGE].changes

Upload Package:
$ dput ppa:michaeldarr/ppa [PACKAGE].changes"