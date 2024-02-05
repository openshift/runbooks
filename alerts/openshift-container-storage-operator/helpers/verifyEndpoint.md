# Verify ODF Provider reachability from ODF Client

You should have the uid of storageconsumer on ODF Provider to which
storageclient on ODF Client established a connection.

Find the provider endpoint configured in ODF Client cluster
```bash
  oc get storageclient -A -ojsonpath='{.items[?(@.status.id==
  "<CONSUMERUID>")].spec.storageProviderEndpoint}'
```

Verify reachability of endpoint gathered from above command, usually
ocs-client-operator is installed in `openshift-storage-client` namespace.
```bash
  oc rsh -n<CLIENT_NAMESPACE> deploy/ocs-client-operator-controller-manager \
  curl <ENDPOINT>
```

Any response other than `Empty reply from server` indicates connection failure.
