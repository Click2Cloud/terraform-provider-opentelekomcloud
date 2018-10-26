#!/bin/sh
#cp to copy build code to folder containing config file


echo "==> Terraform initialization in process..."

sudo cp terraform-provider-opentelekomcloud /positive_testscripts

terraform init ./positive_testscripts

echo "==> Preparing for terraform apply..."

terraform apply -auto-approve ./positive_testscripts

echo "==> Resource destroy in process..."

terraform destroy -force ./positive_testscripts