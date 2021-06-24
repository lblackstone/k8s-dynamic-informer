package informer

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func Watch(inf informers.GenericInformer, runner func(<-chan struct{}, cache.SharedIndexInformer)) {
	stopper := make(chan struct{})
	defer close(stopper)

	runner(stopper, inf.Informer())
}

func PodLogger(stopCh <-chan struct{}, s cache.SharedIndexInformer) {
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
