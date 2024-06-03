#!/bin/bash

# Define the repository details
REPO_NAME="logfire-cli"
REPO_URL="https://logfire-ai.github.io/cli/yum-repo/"
GPG_KEY_URL="https://logfire-ai.github.io/cli/KEY.gpg"

# Create the .repo file
echo "Creating $REPO_NAME.repo in /etc/yum.repos.d/..."
tee /etc/yum.repos.d/$REPO_NAME.repo <<EOL
[$REPO_NAME]
name=Logfire CLI Repository
baseurl=$REPO_URL
enabled=1
gpgcheck=0
gpgkey=$GPG_KEY_URL
EOL

# Import the GPG key
echo "Importing GPG key..."
rpm --import $GPG_KEY_URL

# Refresh the YUM cache
echo "Refreshing YUM cache..."
yum makecache

echo "Logfire CLI Repository has been added and YUM cache has been refreshed."
