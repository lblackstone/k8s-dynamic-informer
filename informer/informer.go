package informer

import (
	"flag"
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func getDynamicInformer(resourceType string) (informers.GenericInformer, error) {
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
		return nil, err
	}
	// Create a factory object that can generate informers for resource types
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dc, 0, corev1.NamespaceAll, nil)
	// "GroupVersionResource" to say what to watch e.g. "deployments.v1.apps" or "pods.v1"
	gvr, gr := schema.ParseResourceArg(resourceType)
	if gvr == nil {
		gvr = &schema.GroupVersionResource{Version: gr.Group, Resource: gr.Resource}
	}
	// Finally, create our informer!
	informer := factory.ForResource(*gvr)
	return informer, nil
}

func PodWatcher(namespace string) {
	//dynamic informer needs to be told which type to watch
	podInformer, _ := getDynamicInformer("pods.v1")
	stopper := make(chan struct{})
	defer close(stopper)
	runPodInformer(stopper, podInformer.Informer(), namespace)
}

func runPodInformer(stopCh <-chan struct{}, s cache.SharedIndexInformer, namespace string) {
	toPod := func(obj interface{}) *corev1.Pod {
		d := &corev1.Pod{}
		err := runtime.DefaultUnstructuredConverter.
			FromUnstructured(obj.(*unstructured.Unstructured).UnstructuredContent(), d)
		if err != nil {
			fmt.Println("could not convert obj to Pod")
			fmt.Print(err)
			return &corev1.Pod{}
		}
		return d
	}
	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := toPod(obj)
			ns := pod.Namespace
			if ns == "" {
				ns = "default"
			}
			fmt.Printf("Added Pod %s/%s\n", ns, pod.Name)
			// do what we want with the Pod event
		},
		DeleteFunc: func(obj interface{}) {
			pod := toPod(obj)
			ns := pod.Namespace
			if ns == "" {
				ns = "default"
			}
			fmt.Printf("Deleted Pod %s/%s\n", ns, pod.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			pod := toPod(newObj)
			ns := pod.Namespace
			if ns == "" {
				ns = "default"
			}
			fmt.Printf("Updated Pod %s/%s\n", ns, pod.Name)
		},
	}
	s.AddEventHandler(handlers)
	s.Run(stopCh)
}
