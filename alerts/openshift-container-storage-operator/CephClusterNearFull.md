# CephClusterNearFull

## Meaning

Storage cluster utilization has crossed 75% and will become read-only at 85%.
Free up some space or expand the storage cluster.

## Impact

Storage cluster will become read-only at 85%.

## Diagnosis

Using the Openshift console, go to Storage-Data Fountation-Storage systems.
A list of the available storage systems with basic information about raw
capacity and used capacity will be visible.
The command "ceph health" provides also information about cluster storage
capacity.

## Mitigation

Two options:

- Scale storage: Depending on the type of cluster it will be needed to add
storage devices and/or nodes. Review the Openshift Scaling storage
documentation.

- Delete information:
If not is possible to scale up the cluster it will be needed to delete
information in order to free space.
