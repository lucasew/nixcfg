FROM mcr.microsoft.com/vscode/devcontainers/universal:linux@sha256:beb3274e5487f10de92c5e02923abe13bd2f2a647e712b32c303a9e91ab7c273

USER root

# [Option] Install Docker CLI
ARG INSTALL_DOCKER="true"
COPY library-scripts/docker-debian.sh /tmp/library-scripts/
RUN if [ "${INSTALL_DOCKER}" = "true" ]; then \
        rm -f /usr/local/share/docker-init.sh \
        && bash /tmp/library-scripts/docker-debian.sh "true" "/var/run/docker-host.sock" "/var/run/docker.sock" "codespace"; \
    fi \
    && rm -rf /var/lib/apt/lists/* /tmp/library-scripts/

USER codespace

# ** [Optional] Uncomment this section to install additional packages. **
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>

RUN curl -L https://nixos.org/nix/install | sh