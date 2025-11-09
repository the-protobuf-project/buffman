FROM ubuntu:22.04

ARG USERNAME=developer
ARG USER_UID=1000
ARG USER_GID=1000
ARG FLATC_VERSION=25.2.10  
ENV DEBIAN_FRONTEND=noninteractive

RUN ln -fs /usr/share/zoneinfo/Asia/Kolkata /etc/localtime

RUN apt-get update && apt-get install -y \
    unzip \
    locales \
    git \
    curl \
    sudo \
    wget \
    patch \
    file \
    ssh 

# Create user 
RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -ms /bin/bash $USERNAME \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME


# Install flatc (remove sudo commands)
RUN wget -O /tmp/flatc.zip https://github.com/google/flatbuffers/releases/download/v${FLATC_VERSION}/Linux.flatc.binary.g++-13.zip && \
    unzip /tmp/flatc.zip -d /tmp && \
    mv /tmp/flatc /usr/local/bin/flatc && \
    chmod +x /usr/local/bin/flatc && \
    rm -rf /tmp/flatc.zip

USER developer