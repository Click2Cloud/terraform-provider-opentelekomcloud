#!/bin/sh
#cp to copy build code to folder containing config file

echo "==> Terraform initialization in process..."

sudo cp terraform-provider-opentelekomcloud /negative_testscripts
terraform init ./negative_testscripts

echo "==> Preparing for terraform apply..."

terraform apply -auto-approve ./negative_testscripts

echo "==> Resource destroy in process..."

terraform destroy -force ./negative_testscripts