# CDIDataImportCronOutdated

## Meaning

This alert fires when `DataImportCron` cannot poll or import the latest disk
image versions.

`DataImportCron` polls disk images, checking for the latest versions, and
imports the images as persistent volume claims (PVCs). This process ensures
that PVCs are updated to the latest version so that they can be used as
reliable clone sources or golden images for virtual machines (VMs).

For golden images, _latest_ refers to the latest operating system of the
distribution. For other disk images, _latest_ refers to the latest hash of the
image that is available.

## Impact

VMs might be created from outdated disk images.

VMs might fail to start because no source PVC is available for cloning.

## Diagnosis

1. Check the cluster for a default storage class:

   ```bash
   $ oc get sc
   ```

   The output displays the storage classes with `(default)` beside the name
   of the default storage class. You must set a default storage class, either on
   the cluster or in the `DataImportCron` specification, in order for the
   `DataImportCron` to poll and import golden images. If no storage class is
   defined, the DataVolume controller fails to create PVCs and the following
   event is displayed: `DataVolume.storage spec is missing accessMode and no
   storageClass to choose profile`.

2. Obtain the `DataImportCron` namespace and name:

   ```bash
   $ oc get dataimportcron -A -o json | jq -r '.items[] | \
     select(.status.conditions[] | select(.type == "UpToDate" and \
     .status == "False")) | .metadata.namespace + "/" + .metadata.name'
   ```

3. If a default storage class is not defined on the cluster, check the
`DataImportCron` specification for a default storage class:

   ```bash
   $ oc get dataimportcron <dataimportcron> -o yaml | \
     grep -B 5 storageClassName
   ```

   Example output:

   ```yaml
         url: docker://.../cdi-func-test-tinycore
       storage:
         resources:
           requests:
             storage: 5Gi
       storageClassName: rook-ceph-block
   ```

4. Obtain the name of the `DataVolume` associated with the `DataImportCron`
object:

   ```bash
   $ oc -n <namespace> get dataimportcron <dataimportcron> -o json | \
     jq .status.lastImportedPVC.name
   ```

5. Check the `DataVolume` log for error messages:

   ```bash
   $ oc -n <namespace> get dv <datavolume> -o yaml
   ```

6. Set the `CDI_NAMESPACE` environment variable:

   ```bash
   $ export CDI_NAMESPACE="$(oc get deployment -A | \
     grep cdi-operator | awk '{print $1}')"
   ```

7. Check the `cdi-deployment` log for error messages:

   ```bash
   $ oc logs -n $CDI_NAMESPACE deployment/cdi-deployment
   ```

## Mitigation

1. Set a default storage class, either on the cluster or in the `DataImportCron`
specification, to poll and import golden images. The updated Containerized Data
Importer (CDI) will resolve the issue within a few seconds.
2. If the issue does not resolve itself, delete the data volumes associated
with the affected `DataImportCron` objects. The CDI will recreate the data
volumes with the default storage class.
3. If your cluster is installed in a restricted network environment, disable
the `enableCommonBootImageImport` feature gate in order to opt out of automatic
updates:

   ```bash
   $ oc patch hco kubevirt-hyperconverged -n $CDI_NAMESPACE --type json \
     -p '[{"op": "replace", "path": \
     "/spec/featureGates/enableCommonBootImageImport", "value": false}]'
   ```

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case, attaching
the artifacts gathered during the diagnosis procedure.