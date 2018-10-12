#!/bin/sh

# This script is used by the Travis build to install terraform packages
sudo apt-get install -y zip unzip
wget https://releases.hashicorp.com/terraform/0.11.7/terraform_0.11.7_linux_amd64.zip
unzip terraform_0.11.7_linux_amd64.zip
sudo mv terraform /usr/local/bin/
terraform --version