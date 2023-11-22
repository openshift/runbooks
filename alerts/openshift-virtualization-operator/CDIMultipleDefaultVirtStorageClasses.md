# CDIMultipleDefaultVirtStorageClasses

## Meaning

This alert fires when more than one storage class has the annotation
`storageclass.kubevirt.io/is-default-virt-class: "true"`.

## Impact

The `storageclass.kubevirt.io/is-default-virt-class: "true"` annotation
defines a default OpenShift Virtualization storage class.

If more than one default OpenShift Virtualization storage class
is defined, a data volume with no storage class specified
receives the most recently created default storage class.

## Diagnosis

Obtain a list of default OpenShift Virtualization storage classes by running
the following command:

```bash
$ oc get sc -o jsonpath='{.items[?(@.metadata.annotations.storageclass\.kubevirt\.io/is-default-virt-class=="true")].metadata.name}'
```

## Mitigation

Ensure that only one default OpenShift Virtualization storage class
is defined by removing the annotation from the other storage classes.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case, attaching
the artifacts gathered during the diagnosis procedure.