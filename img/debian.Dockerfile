FROM ubuntu:18.04

ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8
ENV DEBIAN_FRONTEND=noninteractive 

RUN apt-get update && apt-get install -q -y --no-install-recommends \
        locales \
        sudo \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd -g 999 dfg \
    && useradd -m -r -u 999 -g dfg dfg \
    && mkdir -p /home/dfg/.config \
    && mkdir -p /home/dfg/.cache \
    && echo en_US.UTF-8 UTF-8 > /etc/locale.gen && locale-gen \
    && echo 'root ALL=(ALL) ALL\ndfg ALL=(ALL) NOPASSWD: ALL\nDefaults env_reset\nDefaults secure_path="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"' >> /etc/sudoers \
    && chmod 440 /etc/sudoers \
    && chown -R dfg:dfg /home/dfg \
    && chmod -R g+rw /home/dfg

SHELL ["/bin/bash", "-c"]

USER dfg