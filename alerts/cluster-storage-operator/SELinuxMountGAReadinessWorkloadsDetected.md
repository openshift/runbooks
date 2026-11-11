# SELinuxMountGAReadinessWorkloadsDetected

## Meaning

There are Pods in the cluster that share the same volume that supports SELinux
mounts and the Pods have incompatible SELinux labels.

In OpenShift 5.1, volumes that support it may be mounted with mount option
`-o context=<SELinux label>`, which significantly reduces time to apply the
label to the volume. On the other hand, all Pods that share the same volume
must use the same SELinux label.

The following volumes support mounting with `-o context`:
* In-tree `iSCSI` and `FC`.
* All CSI drivers that explicitly announce support in their `CSIDriver` instance.
  * CSI drivers that are shipped by OpenShift: `AWS EBS`, `Azure Disk`,
    `GCE PersistentDisk`, `OpenStack Cinder`, `vSphere disk`
  * Check `CSIDriver` instance of 3rd party CSI drivers. `seLinuxMount` field
    is `true` for drivers that support SELinux mount option.

All other volumes do not support `-o context` mount option and SELinux is
applied recursively to every single file of Pod volume when the Pod starts.
For a large number of files, it can take a significant amount of time.

## Impact

All Pods that use a volume that supports SELinux mount **must** have the same
SELinux label and the same `spec.securityContext.seLinuxChangePolicy`.

In OpenShift 4.22 / 4.23 / 5.0: All these Pods _may_ run. But these conflicts
must be resolved before upgrade to 5.1, see below.

In OpenShift 5.1 or newer: **Pods with a different SELinux label or
`seLinuxChangePolicy` will not run and may get stuck at `ContainerCreating`.**

## Diagnosis

Find all affected Pods by querying a metric in Console → Observe → Metrics:

```promql
selinux_warning_controller_selinux_volume_conflict > 0
```

The metric has an entry for each conflicting pod pair. Examples:

* ```text
    selinux_warning_controller_selinux_volume_conflict{
        pod1_name="myapp-0",
        pod1_namespace="myapp",
        pod1_value=":::s0:c28,c17",
        pod2_name="my-privileged-pod",
        pod2_namespace="myapp",
        pod2_value="",
        property="SELinuxLabel"} 1
    ```

  * Label `property` shows what field is conflicting - in this case it is
    the SELinux label.
  * `pod1_namespace` and `pod1_name` identify the first pod,
    `pod1_value` shows the conflicting value (SELinux label, in this case).
    The SELinux label has format `user:role:type:level`, where any part may
    be empty. In this example, only the `level` part `s0:c28,c17` is set.
  * `pod2_namespace`, `pod2_name` and `pod2_value` identify the second
    pod and its SELinux label.
  * The metric value is always 1.

  Rewording the above, the metric value says that Pod `myapp/myapp-0` has
  SELinux label `s0:c28,c17` and it shares a volume with Pod
  `myapp/my-privileged-pod`, which has an empty SELinux label (which is
  typical for privileged pods). Sharing a volume among Pods with different
  SELinux labels will not be possible.

* ```text
    selinux_warning_controller_selinux_volume_conflict{
        pod1_name="testpod-recursive",
        pod1_namespace="myapp",
        pod1_value="Recursive",
        pod2_name="myapp-0",
        pod2_namespace="myapp",
        pod2_value="MountOption",
        property="SELinuxChangePolicy"} 1
    ```

    It shows that two Pods have the same SELinux label, but they
    have a different `securityContext.seLinuxChangePolicy` set, which is invalid
    too.

    All Pods that access the same volume concurrently must have the same
    `seLinuxChangePolicy`. If the field is empty, `MountOption` is implied.

### HyperShift Hosted Control Plane

The metric that lists problematic Pods is emitted by kube-controller-manager.

If the cluster control plane runs as a hosted control plane in HyperShift,
then the cluster admin must enable collection of the hosted control plane
metrics. See
[OCP docs](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/hosted_control_planes/hcp-observability)
for how to enable observability for hosted control planes.

## Mitigation

There are several ways to solve this issue:

1. Fix all Pods that share the same volume to have the same SELinux label.
   This may be the best option for unprivileged pods. Typically, edit your
   Deployment / StatefulSet resources and ensure Pod
   `securityContext.seLinuxOptions` is the same and they run with the same SCC.

2. Opt out Pods from the SELinux mount feature and let the container runtime relabel
   all files on the volume when a Pod starts. This allows a privileged and
   unprivileged Pod to share data on a volume concurrently. Still, all unprivileged
   Pods must have the same SELinux label, otherwise the Linux kernel will
   deny any access to the volume to some of these Pods.

   This can be done by setting `securityContext.seLinuxChangePolicy` to
   `Recursive` of all Pods that need to access the volume.

   ```sh
   # oc patch statefulsets/myapp --type=merge -p '{"spec": {"template": { "spec": { "securityContext": { "seLinuxChangePolicy": "Recursive" }}}}}'
   ```

3. Opt out a whole namespace from the SELinux mount feature and set
   `seLinuxChangePolicy` to `Recursive` for any newly created Pods in the
   namespace by labelling the namespace with
   `storage.openshift.io/selinux-change-policy: Recursive`. For example:

   ```sh
   # oc label namespace/myapp "storage.openshift.io/selinux-change-policy=Recursive"
   ```

   This does not change any existing Pods, only newly created ones!

   It is a quick way to update a namespace whose Pods are created by 3rd
   party tools, operators or Helm charts that are out of your control or
   hard to update.
