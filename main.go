package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/AashishsMohankar/terraform-provider-opentelekomcloud/opentelekomcloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: opentelekomcloud.Provider})
}
