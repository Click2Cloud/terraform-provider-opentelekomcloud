#!/bin/sh
#cp to copy build code to folder containing config file
#rm -rf to delete tf.state file to import

echo "==> Terraform initialization in process..."

sudo cp terraform-provider-opentelekomcloud /modify-resource_testscripts

terraform init ./modify-resource_testscripts

echo "==> Preparing for terraform apply..."

terraform apply -auto-approve ./modify-resource_testscripts

sudo rm -rf terraform.tfstate /modify-resource_testscripts

echo "==> for import a resource"

terraform import "${opentelekomcloud_sfs_file_system_v2.Share_file_011.id}" ./modify-resource_testscripts

echo "==> Preparing for terraform apply..."

terraform apply -auto-approve ./modify-resource_testscripts/modify-sfs.tf

echo "==> Resource destroy in process..."

terraform destroy -force ./modify-resource_testscripts