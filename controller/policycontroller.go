package controller

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	crv1 "github.com/crunchydata/kraken/apis/cr/v1"
)

// Watcher is an example of watching on resource create/update/delete events
type PgpolicyController struct {
	PgpolicyClient *rest.RESTClient
	PgpolicyScheme *runtime.Scheme
}

// Run starts an Example resource controller
func (c *PgpolicyController) Run(ctx context.Context) error {
	fmt.Print("Watch Pgpolicy objects\n")

	// Watch Example objects
	_, err := c.watchPgpolicys(ctx)
	if err != nil {
		fmt.Printf("Failed to register watch for Pgpolicy resource: %v\n", err)
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}

func (c *PgpolicyController) watchPgpolicys(ctx context.Context) (cache.Controller, error) {
	source := cache.NewListWatchFromClient(
		c.PgpolicyClient,
		crv1.PgpolicyResourcePlural,
		apiv1.NamespaceAll,
		fields.Everything())

	_, controller := cache.NewInformer(
		source,

		// The object type.
		&crv1.Pgpolicy{},

		// resyncPeriod
		// Every resyncPeriod, all resources in the cache will retrigger events.
		// Set to 0 to disable the resync.
		0,

		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.onAdd,
			UpdateFunc: c.onUpdate,
			DeleteFunc: c.onDelete,
		})

	go controller.Run(ctx.Done())
	return controller, nil
}

func (c *PgpolicyController) onAdd(obj interface{}) {
	example := obj.(*crv1.Pgpolicy)
	fmt.Printf("[PgpolicyCONTROLLER] OnAdd %s\n", example.ObjectMeta.SelfLink)

	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use exampleScheme.Copy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	copyObj, err := c.PgpolicyScheme.Copy(example)
	if err != nil {
		fmt.Printf("ERROR creating a deep copy of example object: %v\n", err)
		return
	}

	exampleCopy := copyObj.(*crv1.Pgpolicy)
	exampleCopy.Status = crv1.PgpolicyStatus{
		State:   crv1.PgpolicyStateProcessed,
		Message: "Successfully processed Pgpolicy by controller",
	}

	err = c.PgpolicyClient.Put().
		Name(example.ObjectMeta.Name).
		Namespace(example.ObjectMeta.Namespace).
		Resource(crv1.PgpolicyResourcePlural).
		Body(exampleCopy).
		Do().
		Error()

	if err != nil {
		fmt.Printf("ERROR updating status: %v\n", err)
	} else {
		fmt.Printf("UPDATED status: %#v\n", exampleCopy)
	}
}

func (c *PgpolicyController) onUpdate(oldObj, newObj interface{}) {
	oldExample := oldObj.(*crv1.Pgpolicy)
	newExample := newObj.(*crv1.Pgpolicy)
	fmt.Printf("[PgpolicyCONTROLLER] OnUpdate oldObj: %s\n", oldExample.ObjectMeta.SelfLink)
	fmt.Printf("[PgpolicyCONTROLLER] OnUpdate newObj: %s\n", newExample.ObjectMeta.SelfLink)
}

func (c *PgpolicyController) onDelete(obj interface{}) {
	example := obj.(*crv1.Pgpolicy)
	fmt.Printf("[PgpolicyCONTROLLER] OnDelete %s\n", example.ObjectMeta.SelfLink)
}
