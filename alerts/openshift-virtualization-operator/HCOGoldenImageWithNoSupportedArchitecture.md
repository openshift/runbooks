# HCOGoldenImageWithNoSupportedArchitecture

## Meaning

When running on a heterogeneous cluster, a cluster with nodes of different
architectures, the DataImportCronTemplates (DICTs; also known as golden
images), in the hyperconverged cluster operator (HCO) should be annotated with
the `ssp.kubevirt.io/dict.architectures` annotation, where the value is the
list of the architectures supported by the image, that is defined in each DICT.

For pre-defined DICTs, this annotation is already set, but for custom DICTs
(user defined DICTs), this annotation must be set by the user in the
HyperConverged custom resource (CR).

For each DICT, if the annotation does not include any architecture that is
supported by the cluster (which mean, there is no node in the cluster with
the architectures listed in the DICT annotation), Then HCO will trigger
the `HCOGoldenImageWithNoSupportedArchitecture` alert for this specific DICT.

> **Note:** This alert is only triggered, if the `enableMultiArchBootImageImport`
> feature gate is enabled in the HyperConverged CR.

## Impact

When this alert is triggered, it means that the DICT is not supported by any of
the nodes in the cluster. HCO will not populate the SSP CR with this DICT, and
so this golden image will not be available for use in the cluster.

## Diagnosis

Read the HyperConverged CR:

```bash
   # Get the namespace of the HyperConverged CR
$ NAMESPACE="$(oc get hyperconverged -A --no-headers | awk '{print $1}')"

#Read the HyperConverged CR
$ oc get hyperconverged -n "${NAMESPACE}" -o yaml
```

There are a few fields in the HyperConverged CR status that can be used to
diagnose this issue:

1. The `status.nodeInfo.workloadsArchitectures` shows the list of architectures
   supported by the cluster.
2. The `status.dataImportCronTemplates` field shows the list of DICTs that are
   managed by HCO.
    1. Find the specific DICT object that is triggering this alert by its name,
       as specified in the alert message. check the DICT's
       `ssp.kubevirt.io/dict.architectures` annotation. Unlike the annotation
       in the spec field, this annotation contain only the architectures that
       are supported by the image **and** by the cluster.

       If the annotation is empty, then there is no architecture supported by
       the image and by the cluster.
    2. The DICT status field will include the `conditions` field, with the
       `Deployed` condition set to `False`, and the `reason` field set to
       `UnsupportedArchitectures`.
       > **Note:** For DICT with supported architectures, the status
         field will not contain the `conditions` field.
    3. The DICT's `status.workloadsArchitectures` field shows the list of
       architectures supported by the image, as was set in the
       `ssp.kubevirt.io/dict.architectures` annotation in the source DICT.

### Example

```yaml
apiVersion: hco.kubevirt.io/v1beta1
kind: HyperConverged
...
status:
  ...
  dataImportCronTemplates:
    - metadata:
        annotations:
          ssp.kubevirt.io/dict.architectures: ""
        name: my-image
      spec:
          ...
      status:
        conditions:
          - message: DataImportCronTemplate has no supported architectures for the current
              cluster
            reason: UnsupportedArchitectures
            status: "False"
            type: Deployed
        originalSupportedArchitectures: someUnsupportedArch,otherUnsupportedArch
```

## Mitigation

### Pre-defined DataImportCronTemplates

The pre-defined DICTs are not defined in the `spec.dataImportCronTemplates`
field in the HyperConverged CR, but they are defined internally in the HCO
application.

All pre-defined DICTs are annotated with the `ssp.kubevirt.io/dict.architectures`
annotation, and all of them supports the `amd64`, `arm64`, and `s390x`
architectures. In the unlikely case that the cluster does not support any of
these architectures, there is no way to use these pre-defined DICTs in the
cluster.

To mitigate this issue, (if adding supported nodes to the cluster is not an
option), you can either:

1. Disable the pre-defined DICTs in the HyperConverged CR, to turn this alert
   off:
    1. Find the DICT(s) you want to disable, in the HyperConverged `status.dataImportCronTemplates`
       field, as described
       [above](#diagnosis).
    2. Add the DICT to the `spec.dataImportCronTemplates` field in the
       HyperConverged CR. Add the `dataimportcrontemplate.kubevirt.io/enable`
       annotation with the value `false` to the DICT. Only the DICT name and
       the annotation are required, in this case

       For example, to disable the `centos-stream10-image-cron` DICT:
       ```yaml
       apiVersion: hco.kubevirt.io/v1beta1
       kind: HyperConverged
       metadata:
         name: kubevirt-hyperconverged
       spec:
         dataImportCronTemplates:
         - metadata:
             name: centos-stream10-image-cron
             annotations:
               dataimportcrontemplate.kubevirt.io/enable: 'false'
         ```
2. If you have the self-built desired image, that is supported by the nodes in
   the cluster, you can modify the pre-defined DICT to use your image, adding
   the DICT to the `spec.dataImportCronTemplates` field in the HyperConverged
   CR, and modify its `spec.source.registry` field.

   > Tip: you can find the pre-defined DICTs in HyperConverged CR `status.dataImportCronTemplates`
   > field, as described [above](#diagnosis). Then you can copy the DICT from
   > there, and modify it in the HyperConverged CR
   > `spec.dataImportCronTemplates` field.

   Don't forget to set the `ssp.kubevirt.io/dict.architectures` annotation to
   include all the architectures supported by your image.

   In this case, you'll need to add all the fields of the DICT.

   For example:
   ```yaml
   apiVersion: hco.kubevirt.io/v1beta1
   kind: HyperConverged
   metadata:
     name: kubevirt-hyperconverged
   spec:
     dataImportCronTemplates:
     - metadata:
       annotations:
         cdi.kubevirt.io/storage.bind.immediate.requested: "true"
         ssp.kubevirt.io/dict.architectures: arch1,arch2
       name: centos-stream10-image-cron
     spec:
       garbageCollect: Outdated
       managedDataSource: centos-stream10
       schedule: "0 */12 * * *"
       template:
         spec:
           source:
             registry:
               url: docker://your-registry/your-image:latest
           storage:
             resources:
               requests:
                 storage: 10Gi
   ```

### User-defined DataImportCronTemplates

User-defined DICTs are defined in the HyperConverged CR, in the
`spec.dataImportCronTemplates` field.

First, check what architectures are supported by the image. You can use the
following command:

```bash
$ podman manifest inspect your-registry/your-image:latest
```

See here for
the [podman manifest inspect
documentation](https://docs.podman.io/en/latest/markdown/podman-manifest-inspect.1.html).

If the image is multi architecture manifest (fat manifest), it will include the
`manifests` field, which is a list of architectures supported by the image. If
the image is not a multi architecture manifest, you will need to find out what
is its architecture.

Then, check that the `ssp.kubevirt.io/dict.architectures` annotation is set
with the correct value. If not, edit the HyperConverged CR to fix the
annotation to the right value. The format of the annotation is a
comma-separated list of architectures; e.g., `amd64,arm64,s390x`.

If the image does not support any of the architectures supported by the
cluster, you will need to either rebuild the image for one or more of
the architectures supported by the cluster, or remove the DICT from the
HyperConverged CR. It is also possible to disable the DICT, by adding it
the `dataimportcrontemplate.kubevirt.io/enable` annotation, with the value
of `false.`; for example:
   ```yaml
   apiVersion: hco.kubevirt.io/v1beta1
   kind: HyperConverged
   metadata:
     name: kubevirt-hyperconverged
   spec:
     dataImportCronTemplates:
     - metadata:
       annotations:
         dataimportcrontemplate.kubevirt.io/enable: "false"
         ssp.kubevirt.io/dict.architectures: unsupported-arch1,unsupported-arch2
       name: my-image
     spec:
       ...
   ```

Find some more information about building multi-architecture images, see the
[podman documentation](https://docs.podman.io/en/latest/markdown/podman-manifest-create.1.html).

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.