# CephFSOrphanedSnapshot

## Meaning

A snapshot is orphaned when there is no attached volume snapshot content available.
This results in data that consumes storage space but remains
inaccessible to the user. This alert is triggered when stale snapshots exist in
a Ceph filesystem.

## Impact

Orphaned snapshots consume storage capacity without providing any value. Over
time, accumulated orphaned snapshots can lead to unexpected storage exhaustion,
potentially impacting the ability to provision new volumes or create new
snapshots.

## Diagnosis

To identify orphaned snapshots in the filesystem, use the ODF CLI tool.

### Install ODF CLI

If not already installed, follow the installation instructions at:
[Odf-Cli Installation](https://github.com/red-hat-storage/odf-cli?tab=readme-ov-file#installation)

### List All Snapshots

To view all snapshots and their binding status:

```bash
odf cephfs-snap ls
```

Snapshots with state `bound` have a corresponding Kubernetes
VolumeSnapshotContent resource. Snapshots with state `orphaned` do not.

### List Only Orphaned Snapshots

To filter and display only the orphaned snapshots:

```bash
odf cephfs-snap ls --orphaned
```

### Running on a Client Cluster

When running the ODF CLI on a client cluster, pass the `--storage-client` flag
with the StorageClient CR name to provide the relevant client cluster details:

```bash
odf cephfs-snap ls --storage-client <storage-client-name>
```

```bash
odf cephfs-snap ls --orphaned --storage-client <storage-client-name>
```

## Mitigation

If the orphaned snapshots are no longer needed and not currently in use, you can
safely delete them to reclaim storage space.

### Delete an Orphaned Snapshot

Use the `delete` subcommand with the subvolume name and snapshot name:

```bash
odf cephfs-snap delete <subvolume> <snapshot>
```

When running on a client cluster, include the `--storage-client` flag:

```bash
odf cephfs-snap delete <subvolume> <snapshot> --storage-client <storage-client-name>
```