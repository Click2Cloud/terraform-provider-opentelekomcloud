#!/bin/sh
#cp to copy build code to folder containing config file

echo "==> Terraform initialization in process..."

sudo cp terraform-provider-opentelekomcloud /query_negative_scenario
terraform init ./negative_testscripts

echo "==> Preparing for terraform apply..."

terraform apply -auto-approve ./query_negative_scenario

echo "==> Resource destroy in process..."

terraform destroy -force ./query_negative_scenario