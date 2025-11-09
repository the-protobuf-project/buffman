#!/bin/bash

# Set the version
BUFFMAN_VERSION="1.0.0"

# Detect operating system
OS=$(uname -s)
case $OS in
    Linux)
        OS_NAME="linux"
        ;;
    Darwin)
        OS_NAME="darwin"
        ;;
    *)
        echo "Unsupported operating system: $OS"
        exit 1
        ;;
esac

# Detect architecture
ARCH=$(uname -m)
case $ARCH in
    x86_64)
        ARCH_NAME="x86-64"
        ;;
    arm64)
        ARCH_NAME="aarch64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Construct download URL
DOWNLOAD_URL="https://github.com/machanirobotics/buffman/releases/download/v$BUFFMAN_VERSION/buffman-$OS_NAME-$ARCH_NAME-$BUFFMAN_VERSION"

echo "Downloading buffman for $OS_NAME-$ARCH_NAME..."
echo "URL: $DOWNLOAD_URL"

# Download the binary
curl -L "$DOWNLOAD_URL" -o buffman

# Check if download was successful
if [ $? -ne 0 ]; then
    echo "Failed to download buffman"
    exit 1
fi

# Move to /usr/local/bin and make executable
sudo mv buffman /usr/local/bin/
sudo chmod +x /usr/local/bin/buffman

echo "buffman installed successfully!"
echo "Version: $BUFFMAN_VERSION"
echo "Platform: $OS_NAME-$ARCH_NAME"
