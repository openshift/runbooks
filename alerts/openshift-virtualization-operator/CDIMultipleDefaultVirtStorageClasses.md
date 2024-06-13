# CDIMultipleDefaultVirtStorageClasses

## Meaning

This alert fires when more than one default virtualization storage class exists.

A default virtualization storage class has precedence over a default OpenShift
Container Platform
storage class for creating a VirtualMachine disk image.

## Impact

If more than one default virtualization storage class exists, a data volume that
requests a default storage class (storage class not explicitly specified),
receives the most recently created one.

## Diagnosis

Obtain a list of default virtualization storage classes by running the following
command:

```bash
$ oc get sc -o json | jq '.items[].metadata|select(.annotations."storageclass.kubevirt.io/is-default-virt-class"=="true")|.name'
```

## Mitigation

Ensure that only one storage class has the default virtualization storage class
annotation.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.