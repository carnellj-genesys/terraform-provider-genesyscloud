package simple_routing_queue

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mypurecloud/platform-client-sdk-go/v105/platformclientv2"
	"log"
	gcloud "terraform-provider-genesyscloud/genesyscloud"
	"terraform-provider-genesyscloud/genesyscloud/consistency_checker"
	"terraform-provider-genesyscloud/genesyscloud/util/resourcedata"
	"time"
)

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
	// Get an instance of the proxy (example can be found in the delete method below)
	sdkConfig := meta.(*gcloud.ProviderMeta).ClientConfig
	proxy := getSimpleRoutingQueueProxy(sdkConfig)

	// Create variables for each field in our schema.ResourceData object
	name := d.Get("name").(string)
	callingPartyName := d.Get("calling_party_name").(string)
	enableTranscription := d.Get("enable_transcription").(bool)

	log.Printf("Creating simple queue %s", name)
	// Create a queue struct using the Genesys Cloud platform go sdk
	queueCreate := &platformclientv2.Createqueuerequest{
		Name:                &name,
		CallingPartyName:    &callingPartyName,
		EnableTranscription: &enableTranscription,
	}

	// Call the proxy function to create our queue
	queueResp, _, err := proxy.createRoutingQueue(ctx, queueCreate)
	if err != nil {
		return diag.Errorf("failed to create queue %s: %v", name, err)
	}

	// Set ID in the schema.ResourceData object
	d.SetId(*queueResp.Id)

	return readSimpleRoutingQueue(ctx, d, meta)
}

// readSimpleRoutingQueue is used by the genesyscloud_simple_routing_queue resource to read a simple queue from Genesys cloud.
func readSimpleRoutingQueue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get an instance of the proxy
	sdkConfig := meta.(*gcloud.ProviderMeta).ClientConfig
	proxy := getSimpleRoutingQueueProxy(sdkConfig)

	log.Printf("Reading simple queue %s", d.Id())
	return gcloud.WithRetriesForRead(ctx, d, func() *resource.RetryError {
		// Call the read queue function to find our queue, passing in the ID from the resource data ( d.Id() )
		// The returned value are: Queue, APIResponse, error
		// If the error is not nil, we should pass the response to the function gcloud.IsStatus404(response)
		// If the status is 404, return a resource.RetryableError. Otherwise, it should be a NonRetryableError
		currentQueue, respCode, err := proxy.getRoutingQueue(ctx, d.Id())
		if err != nil {
			if gcloud.IsStatus404ByInt(respCode) {
				return resource.RetryableError(fmt.Errorf("failed to read queue %s: %v", d.Id(), err))
			}
			return resource.NonRetryableError(fmt.Errorf("failed to read queue %s: %v", d.Id(), err))
		}

		// Define consistency checker
		cc := consistency_checker.NewConsistencyCheck(ctx, d, meta, ResourceSimpleRoutingQueue())

		// Set our values in the schema resource data, based on the values
		// in the Queue object returned from the API
		_ = d.Set("name", *currentQueue.Name)
		resourcedata.SetNillableValue(d, "calling_party_name", currentQueue.CallingPartyName)
		resourcedata.SetNillableValue(d, "enable_transcription", currentQueue.EnableTranscription)

		return cc.CheckState()
	})
}

// updateSimpleRoutingQueue is used by the genesyscloud_simple_routing_queue resource to update a simple queue in Genesys cloud.
func updateSimpleRoutingQueue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get an instance of the proxy
	sdkConfig := meta.(*gcloud.ProviderMeta).ClientConfig
	proxy := getSimpleRoutingQueueProxy(sdkConfig)

	log.Printf("Updating simple queue %s", d.Id())

	// Create variables for each field in our schema.ResourceData object
	name := d.Get("name").(string)
	callingPartyName := d.Get("calling_party_name").(string)
	enableTranscription := d.Get("enable_transcription").(bool)

	// Create a queue struct using the Genesys Cloud platform go sdk
	queueUpdate := &platformclientv2.Queuerequest{
		Name:                &name,
		CallingPartyName:    &callingPartyName,
		EnableTranscription: &enableTranscription,
	}

	// Call the proxy function to update our queue, passing in the queue ID and the queue object
	// All we need from the response is the error for our error handling
	_, _, err := proxy.updateRoutingQueue(ctx, d.Id(), queueUpdate)
	if err != nil {
		return diag.Errorf("failed to update queue %s: %v", name, err)
	}

	return readSimpleRoutingQueue(ctx, d, meta)
}

// deleteSimpleRoutingQueue is used by the genesyscloud_simple_routing_queue resource to delete a simple queue from Genesys cloud.
func deleteSimpleRoutingQueue(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get an instance of the proxy
	sdkConfig := meta.(*gcloud.ProviderMeta).ClientConfig
	proxy := getSimpleRoutingQueueProxy(sdkConfig)

	log.Printf("Deleting simple queue %s", d.Id())

	// Call the delete queue proxy function, passing in our queue ID from the schema.ResourceData object
	_, err := proxy.deleteRoutingQueue(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete queue %s: %v", d.Id(), err)
	}

	// Check that queue has been deleted by trying to get it from the API
	time.Sleep(5 * time.Second)
	return gcloud.WithRetries(ctx, 30*time.Second, func() *resource.RetryError {
		_, respCode, err := proxy.getRoutingQueue(ctx, d.Id())

		if err == nil {
			return resource.NonRetryableError(fmt.Errorf("error deleting routing queue %s: %s", d.Id(), err))
		}
		if gcloud.IsStatus404ByInt(respCode) {
			// Success: Routing Queue deleted
			log.Printf("Deleted routing queue %s", d.Id())
			return nil
		}

		return resource.RetryableError(fmt.Errorf("routing queue %s still exists", d.Id()))
	})
}
