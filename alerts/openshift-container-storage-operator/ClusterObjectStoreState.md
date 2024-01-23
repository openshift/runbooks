# ClusterObjectStoreState

## Meaning

RGW endpoint of the Ceph object store is in a failure state,
OR
One or more Rook Ceph RGW deployments have fewer ready replicas than required
for more than 15s.

## Impact

Cluster Object Store is in unhealthy state
OR
Number of ready replicas for Rook Ceph RGW deployments is less than the desired replicas.

## Diagnosis

Need to check whether the given RGW endpoints are accessible or not.
Make sure that the Ceph RGW deployments have required number of replicas.

## Mitigation

Please check the health of the Ceph cluster and the deployments and find the
root cause of the issue.

