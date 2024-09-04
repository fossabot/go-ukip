#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Ensure the script is run as root
if [ "$(id -u)" != "0" ]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

# Create necessary directories
mkdir -p /etc/ukip
mkdir -p /usr/local/bin

# Copy the main binary
cp cmd/ukip/ukip /usr/local/bin/ukip
chmod +x /usr/local/bin/ukip

# Copy configuration files
cp configs/allowlist.txt /etc/ukip/allowlist
cp configs/keycodes.json /etc/ukip/keycodes

# Copy and enable the systemd service
cp configs/ukip.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable ukip.service
systemctl start ukip.service

echo "UKIP has been installed and started successfully."