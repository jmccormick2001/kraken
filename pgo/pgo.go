package main

import (
	//"context"
	"flag"
	"fmt"
	"time"

	apiv1 "k8s.io/api/core/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	crv1 "github.com/crunchydata/foo/apis/cr/v1"
	exampleclient "github.com/crunchydata/foo/client"
	//examplecontroller "github.com/crunchydata/foo/controller"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "", "Path to a kube config. Only required if out-of-cluster.")
	flag.Parse()

	// Create the client config. Use kubeconfig if given, otherwise assume in-cluster.
	config, err := buildConfig(*kubeconfig)
	if err != nil {
		panic(err)
	}

	//apiextensionsclientset, err := apiextensionsclient.NewForConfig(config)
	_, err = apiextensionsclient.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// make a new config for our extension's API group, using the first config as a baseline
	//exampleClient, exampleScheme, err := exampleclient.NewClient(config)
	exampleClient, _, err := exampleclient.NewClient(config)
	if err != nil {
		panic(err)
	}

	// Create an instance of our custom resource
	example := &crv1.Example{
		ObjectMeta: metav1.ObjectMeta{
			Name: "example1",
		},
		Spec: crv1.ExampleSpec{
			Foo: "hello",
			Bar: true,
		},
		Status: crv1.ExampleStatus{
			State:   crv1.ExampleStateCreated,
			Message: "Created, not processed yet",
		},
	}
	var result crv1.Example
	err = exampleClient.Post().
		Resource(crv1.ExampleResourcePlural).
		Namespace(apiv1.NamespaceDefault).
		Body(example).
		Do().Into(&result)
	if err == nil {
		fmt.Printf("CREATED: %#v\n", result)
	} else if apierrors.IsAlreadyExists(err) {
		fmt.Printf("ALREADY EXISTS: %#v\n", result)
	} else {
		panic(err)
	}

	userLabels := make(map[string]string)
	masterStorage := crv1.PgStorageSpec{}
	masterStorage.PvcName = ""
	masterStorage.StorageClass = ""
	masterStorage.PvcAccessMode = ""
	masterStorage.PvcSize = ""
	masterStorage.StorageType = ""
	masterStorage.FSGROUP = ""
	masterStorage.SUPPLEMENTAL_GROUPS = ""

	clusterexample := &crv1.PgCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster1",
		},
		Spec: crv1.PgClusterSpec{
			Name:                  "cluster1",
			ClusterName:           "cluster1",
			Policies:              "",
			CCP_IMAGE_TAG:         "",
			POSTGRES_FULL_VERSION: "",
			Port:                 "",
			NodeName:             "",
			MasterStorage:        masterStorage,
			ReplicaStorage:       masterStorage,
			PG_MASTER_HOST:       "",
			PG_MASTER_USER:       "",
			PG_MASTER_PASSWORD:   "",
			PG_USER:              "",
			PG_PASSWORD:          "",
			PG_DATABASE:          "",
			PG_ROOT_PASSWORD:     "",
			REPLICAS:             "",
			STRATEGY:             "",
			SECRET_FROM:          "",
			BACKUP_PVC_NAME:      "",
			BACKUP_PATH:          "",
			PGUSER_SECRET_NAME:   "",
			PGROOT_SECRET_NAME:   "",
			PGMASTER_SECRET_NAME: "",
			STATUS:               "",
			PSW_LAST_UPDATE:      "",
			UserLabels:           userLabels,
		},
		Status: crv1.PgClusterStatus{
			State:   crv1.PgClusterStateCreated,
			Message: "Created, not processed yet",
		},
	}
	var clusterresult crv1.PgCluster
	err = exampleClient.Post().
		Resource(crv1.PgClusterResourcePlural).
		Namespace(apiv1.NamespaceDefault).
		Body(clusterexample).
		Do().Into(&clusterresult)
	if err == nil {
		fmt.Printf("CREATED: %#v\n", result)
	} else if apierrors.IsAlreadyExists(err) {
		fmt.Printf("ALREADY EXISTS: %#v\n", result)
	} else {
		fmt.Println(err.Error())
		panic(err)
	}

	// Poll until Example object is handled by controller and gets status updated to "Processed"
	err = exampleclient.WaitForExampleInstanceProcessed(exampleClient, "example1")
	if err != nil {
		panic(err)
	}
	fmt.Print("PROCESSED sleeping 20 seconds\n")

	time.Sleep(20000 * time.Millisecond)

	// Fetch a list of our TPRs
	exampleList := crv1.ExampleList{}
	err = exampleClient.Get().Resource(crv1.ExampleResourcePlural).Do().Into(&exampleList)
	if err != nil {
		panic(err)
	}
	fmt.Printf("LIST: %#v\n", exampleList)
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}
