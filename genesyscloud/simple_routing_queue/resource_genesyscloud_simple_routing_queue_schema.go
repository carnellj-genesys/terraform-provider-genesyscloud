package simple_routing_queue

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	registrar "terraform-provider-genesyscloud/genesyscloud/resource_register"
	// gcloud "terraform-provider-genesyscloud/genesyscloud"
)

const resourceName = "genesyscloud_simple_routing_queue"

func SetRegistrar(l registrar.Registrar) {
	l.RegisterResource(resourceName, ResourceSimpleRoutingQueue())
	l.RegisterDataSource(resourceName, DataSourceSimpleRoutingQueue())
}

func ResourceSimpleRoutingQueue() *schema.Resource {
	return &schema.Resource{
		Description: "Genesys Cloud Simple Routing Queue",

		// CREATE-TODO: Specify our our functions that we defined in resource_genesyscloud_simple_routing_queue.go for performing CRUD operations.
		// For example:
		// ReadContext: gcloud.ReadWithPooledClient(readSimpleRoutingQueue)
		// (The line to import the gcloud package is commented out above - you may have to uncomment this line)

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		Schema:        map[string]*schema.Schema{
			/*
				CREATE-TODO: Define the following three fields:
				1. "name"                 | type: string  | required | description: "The name of our routing queue."
				2. "calling_party_name"   | type: string  | optional | description: "The name to use for caller identification for outbound calls from this queue."
				3. "enable_transcription" | type: boolean | optional | description: "Indicates whether voice transcription is enabled for this queue."
				Note: The field "enable_transcription" is also Computed. This lets the provider know that the API will compute and return a value, should
					the user not specify one in the resource config.
			*/
		},
	}
}

func DataSourceSimpleRoutingQueue() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for Genesys Cloud Simple Routing Queues.",
		// CREATE-TODO: As above, specify the function dataSourceSimpleRoutingQueueRead as the ReadContext of this Resource object

		Schema: map[string]*schema.Schema{
			/*
				CREATE-TODO: Define the only field in our data source:
				"name" | type: string | required | description: "The name of our routing queue."
			*/
		},
	}
}
