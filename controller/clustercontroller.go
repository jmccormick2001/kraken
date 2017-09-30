package controller

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	crv1 "github.com/crunchydata/foo/apis/cr/v1"
)

// Watcher is an example of watching on resource create/update/delete events
type PgClusterController struct {
	PgClusterClient *rest.RESTClient
	PgClusterScheme *runtime.Scheme
}

// Run starts an Example resource controller
func (c *PgClusterController) Run(ctx context.Context) error {
	fmt.Print("Watch PgCluster objects\n")

	// Watch Example objects
	_, err := c.watchPgClusters(ctx)
	if err != nil {
		fmt.Printf("Failed to register watch for PgCluster resource: %v\n", err)
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}

func (c *PgClusterController) watchPgClusters(ctx context.Context) (cache.Controller, error) {
	source := cache.NewListWatchFromClient(
		c.PgClusterClient,
		crv1.PgClusterResourcePlural,
		apiv1.NamespaceAll,
		fields.Everything())

	_, controller := cache.NewInformer(
		source,

		// The object type.
		&crv1.PgCluster{},

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

func (c *PgClusterController) onAdd(obj interface{}) {
	example := obj.(*crv1.PgCluster)
	fmt.Printf("[PgClusterCONTROLLER] OnAdd %s\n", example.ObjectMeta.SelfLink)

	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use exampleScheme.Copy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	copyObj, err := c.PgClusterScheme.Copy(example)
	if err != nil {
		fmt.Printf("ERROR creating a deep copy of example object: %v\n", err)
		return
	}

	exampleCopy := copyObj.(*crv1.PgCluster)
	exampleCopy.Status = crv1.PgClusterStatus{
		State:   crv1.PgClusterStateProcessed,
		Message: "Successfully processed PgCluster by controller",
	}

	err = c.PgClusterClient.Put().
		Name(example.ObjectMeta.Name).
		Namespace(example.ObjectMeta.Namespace).
		Resource(crv1.PgClusterResourcePlural).
		Body(exampleCopy).
		Do().
		Error()

	if err != nil {
		fmt.Printf("ERROR updating status: %v\n", err)
	} else {
		fmt.Printf("UPDATED status: %#v\n", exampleCopy)
	}
}

func (c *PgClusterController) onUpdate(oldObj, newObj interface{}) {
	oldExample := oldObj.(*crv1.PgCluster)
	newExample := newObj.(*crv1.PgCluster)
	fmt.Printf("[PgClusterCONTROLLER] OnUpdate oldObj: %s\n", oldExample.ObjectMeta.SelfLink)
	fmt.Printf("[PgClusterCONTROLLER] OnUpdate newObj: %s\n", newExample.ObjectMeta.SelfLink)
}

func (c *PgClusterController) onDelete(obj interface{}) {
	example := obj.(*crv1.PgCluster)
	fmt.Printf("[PgClusterCONTROLLER] OnDelete %s\n", example.ObjectMeta.SelfLink)
}
