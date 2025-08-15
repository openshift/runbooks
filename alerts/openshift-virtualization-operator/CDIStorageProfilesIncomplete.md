# CDIStorageProfilesIncomplete

## Meaning

This alert fires when a Containerized Data Importer (CDI) storage profile is
incomplete.

If a storage profile is incomplete, the CDI cannot infer persistent volume claim
(PVC) fields, such as `volumeMode` and  `accessModes`, which are required to
create a virtual machine (VM) disk.

## Impact

The CDI cannot create a VM disk on the PVC.

## Diagnosis

- Identify the incomplete storage profile:

  ```bash
  $ oc get storageprofile <storage_class>
  ```

## Mitigation

- Add the missing storage profile information:

  ```bash
  $ oc patch storageprofile local --type=merge -p '{"spec": \
    {"claimPropertySets": [{"accessModes": ["ReadWriteOnce"], \
    "volumeMode": "Filesystem"}]}}'
  ```

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.