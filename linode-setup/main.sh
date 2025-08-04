#!/bin/bash

# TODO

# as ROOT
# [x] setup users
# [x] setup firewall - ssh  
# [] add github ssh so CI/CD can push to linode
# [x] disable password authentication and root login

# as NON_ROOT_USER
# [x] install docker
# [x] install gdal
# [x] install caddy + create Caddyfile
# ---
# [] login to ghcr.io (include token in initial .env)
# ---
# [] sync docker-compose.prod and .env (in the future from github CI/CD)
# [] run docker compose
# -- 
# [] ingest GADM data
# [] run migrations
# [] restart caddy 


# Load environment variables from .env if it exists in the same directory as this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [ -f "$SCRIPT_DIR/.env" ]; then
    export $(grep -v '^#' "$SCRIPT_DIR/.env" | xargs)
else
    echo ".env file not found in $SCRIPT_DIR. Skipping environment variable loading."
    exit 1
fi

# setup by root user. ssh password authentication. 
# call it from ./secrets directory
# TODO: copy files individually  
echo "Creating destination directory on remote server"
ssh $ROOT_USER@$REMOTE_HOST "mkdir -p /tmp/boot_gadm"

echo "Copying files to remote server: $REMOTE_HOST" "as $ROOT_USER"
rsync -avz -e ssh \
   ./setup-users-and-ssh.sh \
   ./.env \
   $ROOT_USER@$REMOTE_HOST:/tmp/boot_gadm

# Check if rsync was successful
if [ $? -ne 0 ]; then
    echo "Error: Failed to copy files to remote server"
    exit 1
fi

# execute initial setup
echo "Executing initial setup on remote server: $REMOTE_HOST" "as $ROOT_USER"
ssh $ROOT_USER@$REMOTE_HOST << 'EOF'
    #!/bin/bash
    set -e # exit on error

    cd /tmp/boot_gadm
    ls -lah

    echo "Check if .env exists"
    if [ ! -f ".env" ]; then
        echo "Error: .env file not found"
        exit 1
    fi

    echo "Check if setup-users-and-ssh.sh exists"
    if [ ! -f "setup-users-and-ssh.sh" ]; then
        echo "Error: setup-users-and-ssh.sh not found"
        exit 1
    fi
    chmod +x setup-users-and-ssh.sh 
    echo "Executing setup-users-and-ssh.sh"
    ./setup-users-and-ssh.sh

    rm /tmp/boot_gadm/setup-users-and-ssh.sh
    rm /tmp/boot_gadm/.env
EOF


SCRIPTS_DIR=/home/$NON_ROOT_USER/scripts

echo "Creating scripts directory on remote server"
ssh -i $SSH_KEY_PATH $NON_ROOT_USER@$REMOTE_HOST "mkdir -p $SCRIPTS_DIR"

echo "Copying files to remote server: $REMOTE_HOST" "as $NON_ROOT_USER"
rsync -avz -e "ssh -i $SSH_KEY_PATH" \
   ./install-docker.sh \
   ./install-caddy.sh \
   ./.env \
   $NON_ROOT_USER@$REMOTE_HOST:$SCRIPTS_DIR 

# execute initial setup
echo "Executing setup on remote server: $REMOTE_HOST" "as $NON_ROOT_USER"
ssh -i $SSH_KEY_PATH $NON_ROOT_USER@$REMOTE_HOST << EOF
    #!/bin/bash

    pwd;
    ls -lah;
    cd $SCRIPTS_DIR
    pwd;

    echo "Check if install-docker.sh exists"
    if [ ! -f "install-docker.sh" ]; then
        echo "Error: install-docker.sh not found"
        exit 1
    fi
    chmod +x install-docker.sh
    echo "Executing install-docker.sh"
    ./install-docker.sh

    echo "--------------------------------"
    echo "Installing GDAL"
    echo "--------------------------------"
    echo $SUDO_PASSWORD | sudo -S apt-get update
    echo $SUDO_PASSWORD | sudo -S apt install -y gdal-bin

    echo "--------------------------------"
    echo "Installing Caddy"
    echo "--------------------------------"
    echo "Check if install-caddy.sh exists"
    if [ ! -f "install-caddy.sh" ]; then
        echo "Error: install-caddy.sh not found"
        exit 1
    fi
    chmod +x install-caddy.sh
    ./install-caddy.sh

    echo "--------------------------------"
    echo "Login to ghcr.io"
    echo "--------------------------------"
    echo $GHCR_TOKEN | docker login ghcr.io -u $GH_USERNAME --password-stdin
    echo "--------------------------------"
EOF