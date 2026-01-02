# CephFSStaleSubvolume

## Meaning

A subvolume in a Ceph filesystem becomes stale when it loses connection to its
associated PersistentVolume. This results in orphaned data that consumes
storage space but remains inaccessible to the user. This alert is triggered when
stale subvolumes exist in a Ceph filesystem.

## Impact

Stale subvolumes consume storage capacity that cannot be accessed or utilized
by OpenShift. This leads to higher-than-expected storage usage and may prevent
storage operations from completing successfully.

## Diagnosis

To identify stale subvolumes in the filesystem, you need to use the ODF CLI
tool.

### Install ODF CLI

If not already installed, follow the installation instructions at:
[Odf-Cli Installation](https://github.com/red-hat-storage/odf-cli?tab=readme-ov-file#installation)

### List All Subvolumes

To view all subvolumes, including those marked as stale:

```bash
odf subvolume ls
```

### List Only Stale Subvolumes

To filter and display only the stale subvolumes:

```bash
odf subvolume ls --stale
```

## Mitigation

If the stale subvolumes are no longer needed and not currently in use, you can
safely delete them to reclaim storage space.

To delete a stale subvolume:

```bash
odf subvolume delete <filesystem> <subvolume> <subvolumegroup>
```

The subvolume,subvolumegroup and filesystem details can be retrieved from the
subvolume list command.