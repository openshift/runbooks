# CephClusterReadOnly

## Meaning

Storage cluster utilization has crossed 85% and will become read-only now.

## Impact

Storage cluster is read-only now and needs immediate data deletion or
 cluster expansion.

Storage cluster will become read-only at 85%.

## Diagnosis

Using the Openshift console, go to Storage-Data Fountation-Storage systems.
A list of the available storage systems with basic information about raw
capacity and used capacity will be visible.
The command "ceph health" provides also information about cluster storage
capacity.

## Mitigation

Free up some space or expand the storage cluster immediately.

Two options:

- Scale storage: Depending on the type of cluster it will be needed to add
storage devices and/or nodes. Review the Openshift Scaling storage
documentation.

- Delete information:
If not is possible to scale up the cluster it will be needed to delete
information in order to free space.
