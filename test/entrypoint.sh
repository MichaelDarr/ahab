#!/usr/bin/env sh
set -e

# Create directories to be used by ahab test suite
mkdir -p /home/ahab/test/project/src /home/ahab/test/empty /home/ahab/.config/ahab

# Example project config
echo '{
    "ahab": "0.1",
    "image": "golang:1.13.7-buster"
}' >> /home/ahab/test/project/ahab.json

# Example user config
echo '{
    "hideCommands": true
}' >> /home/ahab/.config/ahab/config.json

# Docker command passthrough
exec "$@"
