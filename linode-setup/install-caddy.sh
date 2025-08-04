#!/bin/bash

echo "--------------------------------"
echo "RUN THIS SCRIPT AS USER WITH SUDO PRIVILEGES"
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

echo "Install Caddy"
# https://caddyserver.com/docs/install#debian-ubuntu-raspbian

echo $SUDO_PASSWORD | sudo -S apt install -y debian-keyring debian-archive-keyring apt-transport-https curl
# Download and add GPG key
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' > /tmp/caddy-key.gpg
echo $SUDO_PASSWORD | sudo -S gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg < /tmp/caddy-key.gpg

# Add repository
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' > /tmp/caddy-repo.list
echo $SUDO_PASSWORD | sudo -S cp /tmp/caddy-repo.list /etc/apt/sources.list.d/caddy-stable.list
echo $SUDO_PASSWORD | sudo -S chmod o+r /usr/share/keyrings/caddy-stable-archive-keyring.gpg
echo $SUDO_PASSWORD | sudo -S chmod o+r /etc/apt/sources.list.d/caddy-stable.list
echo $SUDO_PASSWORD | sudo -S apt update
echo $SUDO_PASSWORD | sudo -S apt install -y caddy

echo "--------------------------------"

echo "Create Caddyfile"
cat <<EOF > ~/Caddyfile
:80 { 
    reverse_proxy :8080
}
EOF

echo "Format Caddyfile"
caddy fmt --overwrite ~/Caddyfile

echo "Stop any existing Caddy processes"
echo $SUDO_PASSWORD | sudo -S pkill caddy 2>/dev/null || true
echo $SUDO_PASSWORD | sudo -S systemctl stop caddy 2>/dev/null || true

echo "Start Caddy with custom config in background"
nohup echo $SUDO_PASSWORD | sudo -S caddy run --config ~/Caddyfile > ~/caddy.log 2>&1 &

echo "Wait for Caddy to start"
sleep 3

echo "Check if Caddy is running"
ps aux | grep -v grep | grep caddy || echo "Caddy not found in process list"

echo "Show Caddy logs"
tail -n 10 ~/caddy.log || echo "No log file found"