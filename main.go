package main

import (
	"github.com/AashishsMohankar/terraform-provider-opentelekomcloud/opentelekomcloud"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: opentelekomcloud.Provider})
}
