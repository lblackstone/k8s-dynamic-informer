package main

import (
	"github.com/lblackstone/informer-test/informer"
)

func main() {
	informer.PodWatcher("")
	select {}
}
