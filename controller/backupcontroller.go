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
type PgBackupController struct {
	PgBackupClient *rest.RESTClient
	PgBackupScheme *runtime.Scheme
}

// Run starts an Example resource controller
func (c *PgBackupController) Run(ctx context.Context) error {
	fmt.Print("Watch PgBackup objects\n")

	// Watch Example objects
	_, err := c.watchPgBackups(ctx)
	if err != nil {
		fmt.Printf("Failed to register watch for PgBackup resource: %v\n", err)
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}

func (c *PgBackupController) watchPgBackups(ctx context.Context) (cache.Controller, error) {
	source := cache.NewListWatchFromClient(
		c.PgBackupClient,
		crv1.PgBackupResourcePlural,
		apiv1.NamespaceAll,
		fields.Everything())

	_, controller := cache.NewInformer(
		source,

		// The object type.
		&crv1.PgBackup{},

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

func (c *PgBackupController) onAdd(obj interface{}) {
	example := obj.(*crv1.PgBackup)
	fmt.Printf("[PgBackupCONTROLLER] OnAdd %s\n", example.ObjectMeta.SelfLink)

	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use exampleScheme.Copy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	copyObj, err := c.PgBackupScheme.Copy(example)
	if err != nil {
		fmt.Printf("ERROR creating a deep copy of example object: %v\n", err)
		return
	}

	exampleCopy := copyObj.(*crv1.PgBackup)
	exampleCopy.Status = crv1.PgBackupStatus{
		State:   crv1.PgBackupStateProcessed,
		Message: "Successfully processed PgBackup by controller",
	}

	err = c.PgBackupClient.Put().
		Name(example.ObjectMeta.Name).
		Namespace(example.ObjectMeta.Namespace).
		Resource(crv1.PgBackupResourcePlural).
		Body(exampleCopy).
		Do().
		Error()

	if err != nil {
		fmt.Printf("ERROR updating status: %v\n", err)
	} else {
		fmt.Printf("UPDATED status: %#v\n", exampleCopy)
	}
}

func (c *PgBackupController) onUpdate(oldObj, newObj interface{}) {
	oldExample := oldObj.(*crv1.PgBackup)
	newExample := newObj.(*crv1.PgBackup)
	fmt.Printf("[PgBackupCONTROLLER] OnUpdate oldObj: %s\n", oldExample.ObjectMeta.SelfLink)
	fmt.Printf("[PgBackupCONTROLLER] OnUpdate newObj: %s\n", newExample.ObjectMeta.SelfLink)
}

func (c *PgBackupController) onDelete(obj interface{}) {
	example := obj.(*crv1.PgBackup)
	fmt.Printf("[PgBackupCONTROLLER] OnDelete %s\n", example.ObjectMeta.SelfLink)
}
