package controller

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	crv1 "github.com/crunchydata/kraken/apis/cr/v1"
	clusteroperator "github.com/crunchydata/kraken/operator/cluster"
)

// Watcher is an example of watching on resource create/update/delete events
type PgclusterController struct {
	PgclusterClient    *rest.RESTClient
	PgclusterScheme    *runtime.Scheme
	PgclusterClientset *kubernetes.Clientset
}

// Run starts an Example resource controller
func (c *PgclusterController) Run(ctx context.Context) error {
	fmt.Print("Watch Pgcluster objects\n")

	// Watch Example objects
	_, err := c.watchPgclusters(ctx)
	if err != nil {
		fmt.Printf("Failed to register watch for Pgcluster resource: %v\n", err)
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}

func (c *PgclusterController) watchPgclusters(ctx context.Context) (cache.Controller, error) {
	source := cache.NewListWatchFromClient(
		c.PgclusterClient,
		crv1.PgclusterResourcePlural,
		apiv1.NamespaceAll,
		fields.Everything())

	_, controller := cache.NewInformer(
		source,

		// The object type.
		&crv1.Pgcluster{},

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

func (c *PgclusterController) onAdd(obj interface{}) {
	example := obj.(*crv1.Pgcluster)
	fmt.Printf("[PgclusterCONTROLLER] OnAdd %s\n", example.ObjectMeta.SelfLink)

	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use exampleScheme.Copy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	copyObj, err := c.PgclusterScheme.Copy(example)
	if err != nil {
		fmt.Printf("ERROR creating a deep copy of example object: %v\n", err)
		return
	}

	exampleCopy := copyObj.(*crv1.Pgcluster)
	exampleCopy.Status = crv1.PgclusterStatus{
		State:   crv1.PgclusterStateProcessed,
		Message: "Successfully processed Pgcluster by controller",
	}

	err = c.PgclusterClient.Put().
		Name(example.ObjectMeta.Name).
		Namespace(example.ObjectMeta.Namespace).
		Resource(crv1.PgclusterResourcePlural).
		Body(exampleCopy).
		Do().
		Error()

	if err != nil {
		fmt.Printf("ERROR updating status: %v\n", err)
	} else {
		fmt.Printf("UPDATED status: %#v\n", exampleCopy)
	}

	clusteroperator.AddClusterBase(c.PgclusterClientset, c.PgclusterClient, exampleCopy, example.ObjectMeta.Namespace)
}

func (c *PgclusterController) onUpdate(oldObj, newObj interface{}) {
	oldExample := oldObj.(*crv1.Pgcluster)
	newExample := newObj.(*crv1.Pgcluster)
	fmt.Printf("[PgclusterCONTROLLER] OnUpdate oldObj: %s\n", oldExample.ObjectMeta.SelfLink)
	fmt.Printf("[PgclusterCONTROLLER] OnUpdate newObj: %s\n", newExample.ObjectMeta.SelfLink)
}

func (c *PgclusterController) onDelete(obj interface{}) {
	example := obj.(*crv1.Pgcluster)
	fmt.Printf("[PgclusterCONTROLLER] OnDelete %s\n", example.ObjectMeta.SelfLink)
}
