#!/bin/bash

# Define the repository details
REPO_NAME="logfire-cli"
REPO_URL="https://logfire-sh.github.io/cli/yum-repo/"
GPG_KEY_URL="https://logfire-sh.github.io/cli/KEY.gpg"

# Create the .repo file
echo "Creating $REPO_NAME.repo in /etc/yum.repos.d/..."
sudo tee /etc/yum.repos.d/$REPO_NAME.repo <<EOL
[$REPO_NAME]
name=Logfire CLI Repository
baseurl=$REPO_URL
enabled=1
gpgcheck=1
gpgkey=$GPG_KEY_URL
EOL

# Import the GPG key
echo "Importing GPG key..."
sudo rpm --import $GPG_KEY_URL

# Refresh the YUM cache
echo "Refreshing YUM cache..."
sudo yum makecache

echo "Logfire CLI Repository has been added and YUM cache has been refreshed."
