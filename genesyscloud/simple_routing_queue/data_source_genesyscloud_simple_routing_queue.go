package simple_routing_queue

import (
	"context"
	"log"
	"terraform-provider-genesyscloud/genesyscloud/provider"
	"terraform-provider-genesyscloud/genesyscloud/util"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
   The data_source_genesyscloud_simple_routing_queue.go contains the data source implementation
   for the resource.
*/

// dataSourceSimpleRoutingQueueRead retrieves by search term the id in question
func dataSourceSimpleRoutingQueueRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// CREATE-TODO 1: Replace '_ =' with 'sdkConfig :=' to capture our SDK configuration in a variable
	sdkConfig := meta.(*provider.ProviderMeta).ClientConfig

	// CREATE-TODO 2: Get an instance of our proxy by passing sdkConfig into the method getSimpleRoutingQueueProxy
	proxy := getSimpleRoutingQueueProxy(sdkConfig)

	// CREATE-TODO 3: Grab our queue name from the schema.ResourceData object (this step is already complete)
	name := d.Get("name").(string)

	log.Printf("Finding queue by name '%s'", name)
	return util.WithRetries(ctx, 15*time.Second, func() *retry.RetryError {
		// CREATE-TODO 4: Call to the proxy function proxyInstance.getRoutingQueueIdByName(context.Context, string), passing ctx and our name variable
		// This function returns values in the following order: queueId (string), response (*platformclientv2.APIResponse), err (error), retryable (bool)
		queueId, resp, err, retryable := proxy.getSimpleRoutingQueueIdByName(ctx, name)

		// CREATE-TODO 5: If the error is not nil, and retryable equals false, return a resource.NonRetryableError
		// to let the user know that an error occurred. If retryable is true, return a resource.RetryableError
		if err != nil {
			if !retryable {
				return retry.NonRetryableError(util.BuildWithRetriesApiDiagnosticError(resourceName, err.Error(), resp))
			}
			return retry.RetryableError(util.BuildWithRetriesApiDiagnosticError(resourceName, err.Error(), resp))
		}

		// CREATE-TODO 6: If we made it this far, we can call d.SetId(queueId) and return nil
		d.SetId(queueId)

		return nil
	})
}
