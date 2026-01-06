# CephMdsMissingReplicas

## Meaning

Minimum required replicas for the storage metadata service (MDS) are not available.
This might affect the working of storage cluster.

## Impact

MDS is responsible for file metadata. Degradation of the MDS service can affect
the working of the storage cluster (related to cephfs storage class) and should
be fixed as soon as possible.

## Diagnosis

Make sure we have enough RAM provisioned for MDS Cache. Default is 4GB, but
recomended is minimum 8GB.

## Mitigation

It is highly recomended to distribute MDS daemons across at least two nodes in
the cluster. Otherwise, a hardware failure on a single node may result in the
file system becoming unavailable.

