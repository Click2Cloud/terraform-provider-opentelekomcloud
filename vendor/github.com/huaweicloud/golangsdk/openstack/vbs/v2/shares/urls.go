package shares

import "github.com/huaweicloud/golangsdk"

const resourcePath = "os-vendor-backup-sharing"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath, "detail")
}