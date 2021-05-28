package main

import (
	"github.com/lblackstone/k8s-dynamic-informer/informer"
)

func main() {
	informer.PodWatcher("")
	select {}
}
