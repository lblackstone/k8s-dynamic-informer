# k8s-dynamic-informer
Testing out Dynamic Informers in Kubernetes

## Running

```
go run ./main.go
```

The output will look something like this:

```
Added Pod kube-system/coredns-f9fd979d6-8rbmd
Added Pod kube-system/coredns-f9fd979d6-vm4p4
Added Pod kube-system/etcd-docker-desktop
Added Pod kube-system/kube-apiserver-docker-desktop
Added Pod kube-system/kube-controller-manager-docker-desktop
Added Pod kube-system/kube-proxy-nwjht
Added Pod kube-system/kube-scheduler-docker-desktop
Added Pod kube-system/storage-provisioner
Added Pod kube-system/vpnkit-controller
Added Pod default/bar
Updated Pod default/bar
Updated Pod default/bar
Updated Pod default/bar
Updated Pod default/bar
Updated Pod default/bar
Updated Pod default/bar
Updated Pod default/bar
Updated Pod default/bar
Deleted Pod default/bar
```

Based on the example from https://hackernoon.com/platforms-on-k8s-with-golang-watch-any-crd-0v2o3z1q
