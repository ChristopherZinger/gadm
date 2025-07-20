#!/bin/bash

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [ -f "$SCRIPT_DIR/.env" ]; then
    export $(grep -v '^#' "$SCRIPT_DIR/.env" | xargs)
else
    echo ".env file not found in $SCRIPT_DIR. Exiting."
    exit 1
fi

echo "NON_ROOT_USER: $NON_ROOT_USER"
echo "USER_PASSWORD: $USER_PASSWORD"

echo "SETUP A UBUNTU VIRTUAL MACHINE"
echo "ENSURE YOU RUN THIS SCRIPT AS ROOT USER"
echo "ENSURE DROPLET HAS PUBLIC SSH KEY ALREADY ADDED TO authorized_keys"
echo ""

echo "--------------------------------"
echo "Adding user $NON_ROOT_USER with password from .env file."

# Create user without password first
adduser --gecos "" --disabled-password $NON_ROOT_USER

# Set password using chpasswd (password is piped to stdin)
echo "$NON_ROOT_USER:$NON_ROOT_USER_PASSWORD" | chpasswd

echo "User $NON_ROOT_USER created with password from .env file"

usermod -aG sudo $NON_ROOT_USER

echo "--------------------------------"
echo "Setup SSH for non sudo user"

# Create .ssh directory for $NON_ROOT_USER user
mkdir -p /home/$NON_ROOT_USER/.ssh

# Copy authorized_keys from root to $NON_ROOT_USER user
cp /root/.ssh/authorized_keys /home/$NON_ROOT_USER/.ssh/authorized_keys

# Set correct ownership
chown -R $NON_ROOT_USER:$NON_ROOT_USER /home/$NON_ROOT_USER/.ssh

# Set correct permissions
chmod 700 /home/$NON_ROOT_USER/.ssh
chmod 600 /home/$NON_ROOT_USER/.ssh/authorized_keys

echo "SSH setup complete for user $NON_ROOT_USER"

echo "--------------------------------"
echo "Setup Firewall"

ufw app list
echo "Allow ssh connections" 
ufw allow OpenSSH
ufw allow 80/tcp comment 'Caddy HTTP'
echo "y" | ufw enable
ufw status

# Disable password authentication
sudo sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
sudo sed -i 's/PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config

# Also disable root login for security
sudo sed -i 's/#PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config
sudo sed -i 's/PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config

# Restart SSH service
sudo systemctl restart ssh

echo "COMPLETED INITIAL SETUP"

exit 0




