package main

import (
	"flag"
	"path/filepath"

	"github.com/lblackstone/k8s-dynamic-informer/informer"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	// Grab a dynamic interface that we can create informers from
	dc, err := dynamic.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}
	factory := dynamicinformer.NewDynamicSharedInformerFactory(dc, 0)
	podInformer := factory.ForResource(schema.GroupVersionResource{Version: "v1", Resource: "pods"})
	deploymentInformer := factory.ForResource(schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"})

	stopper := make(chan struct{})
	defer close(stopper)

	factory.Start(stopper)
	go informer.DeploymentLogger(stopper, deploymentInformer.Informer())
	go informer.PodLogger(stopper, podInformer.Informer())

	select {}
}
