# StorageClientHeartbeatMissed

## Meaning

StorageConsumer heartbeat isn't received from connected storage clients.

## Impact

Ceph monitor endpoints at storage client will have stale information. If
monitor endpoints are changed after loosing heartbeat, storage clients may not
be able to connect to ceph monitors.

## Diagnosis

Take a note of storageconsumer name from alert description and find the
connected client cluster by following
[connectedClient](helpers/connectedClient.md) document.

Login to the client cluster identified by the cluster id from above document.

Verify ODF Provider reachability from ODF Client by following
[verifyEndpoint](helpers/verifyEndpoint.md) document.

## Mitigation

### Intermittent nework connectivity

1. From diagnosis if you find endpoint is reachable, wait for at most 5 minutes
to have connection reestablished which should stop firing alert.
2. If the alert is still firing, check for any outage in your network which is
effecting ODF Provider and ODF Client connectivity, specifically the NodePort
on ODF Provider cluster referenced by storageclient resource.

### Wrong endpoint configured

1. Make sure the endpoint referenced by storageclient resource matches the
endpoint in storagecluster status
``` bash
 oc get -nopenshift-storage storagecluster ocs-storagecluster \
 -ojsonpath='{.status.storageProviderEndpoint}'
```
