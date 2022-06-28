# GarbageCollectorSyncFailed

## Meaning

[This alert][GarbageCollectorSyncFailed] is fired when some resources couldn't
be reached for a garbage collector informer cache. This usually happens due to
the failing conversion webhook of an installed `CustomResourceDefinition` (CRD)
or because of failing API Discovery, unreachable API server or problems with a
network.

## Impact

Garbage collection and orphaning of objects stops working. Normal deletion
works fine, but deletion of objects via owner references does not.
Deletion of an object with `--cascade=orphan` option also does not work.

This can lead to resource exhaustion or/and etcd pollution
when objects are not garbage collected.

For example:
- deletion of replica set will not delete its pods and the pods will keep
  running
- cron job with history limit that create new jobs will delete old jobs,
  but it will lead to leaving completed pods behind

In extreme cases this can lead to cluster failure.


## Diagnosis

Analyze logs of `kube-controller-manager` pods in
`openshift-kube-controller-manager` namespace.

```console
$ oc get pods -n openshift-kube-controller-manager
$ oc logs -n openshift-kube-controller-manager $POD
```

Look for lines with `garbagecollector` or `garbage` words.
You should see something similar to this.

```text
shared_informer.go:258] unable to sync caches for garbage collector
garbagecollector.go:245] timed out waiting for dependency graph builder sync during GC sync (attempt 26)
garbagecollector.go:215] syncing garbage collector with updated resources from discovery (attempt 27): added: [example.com/v1, Resource=myresource], removed: []
```

If you cannot identify the failing resource,
increase the log level to get more information.

```console
$ oc patch kubecontrollermanagers.operator/cluster --type=json -p '[{"op": "replace", "path": "/spec/logLevel", "value": "LOGLEVEL" }]'
```

After the kube-controller-manager pods restart,
you should find following messages in logs.

```console
graph_builder.go:279] garbage controller monitor not yet synced: example.com/v1, Resource=myresource
```


## Mitigation

Debug your CRD (in this example `myresource.example.com`) to see what is wrong.
Fix or disable the conversion webhooks of your CRD to get garbage collector
back to functioning state. You can also remove the CRD in case you are not
using it. Usually this is caused by improper CRD installation or upgrade.
For more details see [Versions in CRDs][VersionsCustomResourceDefinitions].

You might also see network errors in kube-controller-manager logs.
In that case there is probably an issue with your infrastructure.

[GarbageCollectorSyncFailed]: https://github.com/openshift/cluster-kube-controller-manager-operator/blob/20179ecfa3b8c5e766a21c98107f45b84196b914/manifests/0000_90_kube-controller-manager-operator_05_alerts.yaml#L42-L50
[VersionsCustomResourceDefinitions]: https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definition-versioning/