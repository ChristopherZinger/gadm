#!/bin/bash

echo "--------------------------------"
echo "Install Docker"
echo "--------------------------------"

# Get the directory where the script is located
echo "SCRIPT_DIR: $SCRIPT_DIR"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [ -f "$SCRIPT_DIR/.env" ]; then
    export $(grep -v '^#' "$SCRIPT_DIR/.env" | xargs)
else
    echo ".env file not found in $SCRIPT_DIR. Exiting."
    exit 1
fi

echo "--------------------------------"
echo "Add Docker's official GPG key"

echo $SUDO_PASSWORD | sudo -S apt-get update
echo $SUDO_PASSWORD | sudo -S apt-get install -y ca-certificates curl
echo $SUDO_PASSWORD | sudo -S install -m 0755 -d /etc/apt/keyrings
echo $SUDO_PASSWORD | sudo -S curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
echo $SUDO_PASSWORD | sudo -S chmod a+r /etc/apt/keyrings/docker.asc

echo "--------------------------------"
echo "Add Docker's repository to Apt sources"
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
echo $SUDO_PASSWORD | sudo -S apt-get update

echo "--------------------------------"
echo "Add Docker's repository to Apt sources"
echo $SUDO_PASSWORD | sudo -S apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

echo "--------------------------------"
echo "Add current user to docker group"
echo $SUDO_PASSWORD | sudo -S usermod -aG docker $NON_ROOT_USER

echo "--------------------------------"
echo "Activate docker group for current session"
newgrp docker

echo $SUDO_PASSWORD | sudo -S docker run hello-world

echo "--------------------------------"
echo "Completed docker installation!"
echo "--------------------------------"