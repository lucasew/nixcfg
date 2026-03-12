FROM mcr.microsoft.com/vscode/devcontainers/universal:linux@sha256:5035562d736642e6d0382d4ef599d84d43f428609ae17fbe2e841b30451de0f7

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