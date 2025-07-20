#!/bin/bash

# TODO

# as ROOT
# [x] setup users
# [x] setup firewall - ssh and caddy
# [x] disable password authentication and root login

# as NON_ROOT_USER
# [-] install caddy
# [-] create caddyfile
# [-] start caddy
# [] install gdal
# [] install postgres and postgis
# [] initialize gadm database
# [] download data from GADM
# [] ingest data into postgres
# [] deploy gadm server [dockerize go server]
# [] run migrations

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
echo "Copying files to remote server: $REMOTE_HOST" "as $ROOT_USER"
rsync -avz -e ssh ./ root@$REMOTE_HOST:/tmp

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

cd /tmp 
ls -lah

# Check if script exists
if [ ! -f "setup-users-and-ssh.sh" ]; then
    echo "Error: setup-users-and-ssh.sh not found"
    exit 1
fi

# Check if .env exists
if [ ! -f ".env" ]; then
    echo "Error: .env file not found"
    exit 1
fi

chmod +x setup-users-and-ssh.sh 

echo "Executing setup-users-and-ssh.sh"
./setup-users-and-ssh.sh
EOF
