package simple_routing_queue

import (
	"context"
	"fmt"
	"log"
	"terraform-provider-genesyscloud/genesyscloud/consistency_checker"
	"terraform-provider-genesyscloud/genesyscloud/provider"
	"terraform-provider-genesyscloud/genesyscloud/util"
	"terraform-provider-genesyscloud/genesyscloud/util/constants"
	"time"

	resourceExporter "terraform-provider-genesyscloud/genesyscloud/resource_exporter"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mypurecloud/platform-client-sdk-go/v133/platformclientv2"
)

func getAllRoutingQueues(ctx context.Context, clientConfig *platformclientv2.Configuration) (resourceExporter.ResourceIDMetaMap, diag.Diagnostics) {
	resources := make(resourceExporter.ResourceIDMetaMap)
	proxy := getSimpleRoutingQueueProxy(clientConfig)

	// Newly created resources often aren't returned unless there's a delay
	time.Sleep(5 * time.Second)

	queues, resp, err := proxy.getAllSimpleRoutingQueues(ctx)
	if err != nil {
		return nil, util.BuildAPIDiagnosticError(resourceName, fmt.Sprintf("failed to get routing queues: %v", err), resp)
	}

	for _, queue := range *queues {
		resources[*queue.Id] = &resourceExporter.ResourceMeta{Name: *queue.Name}
	}

	return resources, nil
}

/*
The resource_genesyscloud_simple_routing_queue.go contains all of the methods that perform the core logic for a resource.
In general a resource should have a approximately 5 methods in it:

1.  A create.... function that the resource will use to create a Genesys Cloud object (e.g. genesyscloud_simple_routing_queue)
2.  A read.... function that looks up a single resource.
3.  An update... function that updates a single resource.
4.  A delete.... function that deletes a single resource.

Two things to note:

1.  All code in these methods should be focused on getting data in and out of Terraform.  All code that is used for interacting
    with a Genesys API should be encapsulated into a proxy class contained within the package.

2.  In general, to keep this file somewhat manageable, if you find yourself with a number of helper functions move them to a
utils file in the package.  This will keep the code manageable and easy to work through.
*/

// createSimpleRoutingQueue is used by the genesyscloud_simple_routing_queue resource to create a simple queue in Genesys cloud.
func createSimpleRoutingQueue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// CREATE-TODO 1: Get an instance of the proxy (example can be found in the delete method below)

	// CREATE-TODO 2: Create variables for each field in our schema.ResourceData object
	// Example:
	// name := d.Get("name").(string)

	log.Printf("Creating simple routing queue")

	// CREATE-TODO 3: Create a queue struct using the Genesys Cloud platform go sdk
	// Here is the source code for the struct we will be using - https://github.com/MyPureCloud/platform-client-sdk-go/blob/master/platformclientv2/createqueuerequest.go
	// (Remember - we only need to worry about the three fields defined in our schema)

	// CREATE-TODO 4: Call the proxy function to create our queue. The proxy function we want to use here is createSimpleRoutingQueue(ctx context.Context, queue *platformclientv2.Createqueuerequest)
	// If the returned error is not nil, use the BuildAPIDiagnosticError function in the util package to build our error message

	// CREATE-TODO 5: Call d.SetId, passing in the ID attached to the queue object that was returned from createSimpleRoutingQueue
	// This will set the ID of this resource in our tf.state file

	log.Println("Created simple routing queue")
	return readSimpleRoutingQueue(ctx, d, meta)
}

// readSimpleRoutingQueue is used by the genesyscloud_simple_routing_queue resource to read a simple queue from Genesys cloud.
func readSimpleRoutingQueue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// CREATE-TODO 1: Get an instance of the proxy

	log.Printf("Reading simple queue %s", d.Id())
	cc := consistency_checker.NewConsistencyCheck(ctx, d, meta, ResourceSimpleRoutingQueue(), constants.DefaultConsistencyChecks, resourceName)
	return util.WithRetriesForRead(ctx, d, func() *retry.RetryError {
		/*
			CREATE-TODO 2: Call the proxy function getSimpleRoutingQueue(ctx context.Context, id string) to find our queue, passing in the ID from the resource data object (can be retrieved
			using the function d.Id())
			The returned value are: *platformclientv2.Queue, *platformclientv2.APIResponse, error
			If the error is not nil, we should pass the returned *platformclientv2.APIResponse object to the function util.IsStatus404(response)
			If the response is a 404, return a retry.RetryableError. Otherwise, it should be a NonRetryableError

			To build our error message, pass the function util.BuildWithRetriesApiDiagnosticError(resourceName string, summary string, response *platformclientv2.APIResponse)
			into the (Non)RetryableError method
		*/

		// CREATE-TODO 3: Set our values in the schema resource data, based on the values in the Queue object returned from the API
		// For fields that are optional, we should check if the returned value is nil before dereferencing it. We can accomplish this using the
		// method SetNillableValue in the resourcedata package. An example of this being done with calling_party_name:
		// resourcedata.SetNillableValue(d, "calling_party_name", queue.CallingPartyName)
		// There is no need to use this for the name field since we know the pointer to the value will never be nil, so we can use the standard d.Set("name", *queue.Name)

		log.Println("Read simple routing queue")
		return cc.CheckState(d)
	})
}

// updateSimpleRoutingQueue is used by the genesyscloud_simple_routing_queue resource to update a simple queue in Genesys cloud.
func updateSimpleRoutingQueue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("Updating simple queue %s", d.Id())
	// CREATE-TODO 1: Get an instance of the proxy

	// CREATE-TODO 2: Create variables for each field in our schema.ResourceData object

	// CREATE-TODO	3: Create a queue struct using the Genesys Cloud platform go sdk
	// ** NOTE: We are using Queuerequest this time - not Createqueuerequest **
	// Here is the source code for the struct we will be using - https://github.com/MyPureCloud/platform-client-sdk-go/blob/master/platformclientv2/Queuerequest.go
	// (Remember - we only need to worry about the three fields defined in our schema)

	// CREATE-TODO 4: Call the proxy function updateSimpleRoutingQueue(context.Context, id string, *platformclientv2.Queuerequest) to update our queue
	// We should handle our error and response objects the same way as in the createSimpleRoutingQueue method above.
	// We won't be needing the returned Queue object, so an underscore can go in that variables place. If we were to define it, Go would complain that we're not using it,
	// so this is our way of telling Go that we don't need it.
	log.Println("Updating simple routing queue")

	log.Println("Updated simple routing queue")
	return readSimpleRoutingQueue(ctx, d, meta)
}

// deleteSimpleRoutingQueue is used by the genesyscloud_simple_routing_queue resource to delete a simple queue from Genesys cloud.
func deleteSimpleRoutingQueue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// CREATE-TODO 1: Get an instance of the proxy (done)
	sdkConfig := meta.(*provider.ProviderMeta).ClientConfig
	proxy := getSimpleRoutingQueueProxy(sdkConfig)

	log.Printf("Deleting simple queue %s", d.Id())

	// CREATE-TODO 2: Call the proxy function deleteSimpleRoutingQueue(ctx context.Context, id string)
	// If the error is not nil, we should handle it as in the create and update functions

	// Check that queue has been deleted by trying to get it from the API
	return util.WithRetries(ctx, 30*time.Second, func() *retry.RetryError {
		_, resp, err := proxy.getSimpleRoutingQueue(ctx, d.Id())
		if err != nil {
			if util.IsStatus404(resp) {
				log.Println("Successfully deleted simple routing queue.")
				return nil
			}
			errorSummary := fmt.Sprintf("unexpected error encountered reading queue %s: %s", d.Id(), err)
			return retry.NonRetryableError(util.BuildWithRetriesApiDiagnosticError(resourceName, errorSummary, resp))
		}
		return retry.RetryableError(util.BuildWithRetriesApiDiagnosticError(resourceName, "Simple routing queue still exists", resp))
	})
}
