# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.177.0/containers/codespaces-linux/.devcontainer/base.Dockerfile

FROM mcr.microsoft.com/vscode/devcontainers/universal:1-focal

# ** [Optional] Uncomment this section to install additional packages. **
USER root
#
RUN add-apt-repository ppa:cncf-buildpacks/pack-cli
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends pack-cli

RUN wget https://github.com/wagoodman/dive/releases/download/v0.9.2/dive_0.9.2_linux_amd64.deb
RUN sudo apt install ./dive_0.9.2_linux_amd64.deb && rm /dive_0.9.2_linux_amd64.deb
#
USER codespace