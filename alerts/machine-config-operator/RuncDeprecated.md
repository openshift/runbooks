# RuncDeprecated

## Meaning

This informational alert fires when the cluster is using the deprecated runc
container runtime instead of the recommended crun runtime.

This alert will fire when:

- One or more nodes in the cluster are running containers using the runc OCI runtime
- The condition persists for 10 minutes or more

## Impact

The runc container runtime has been deprecated and support will be removed
in a future OpenShift release. While runc continues to function in the
current release, clusters using runc will need to migrate to crun before
upgrading to a future version where runc support is removed.

Continuing to use runc may:
- Block future cluster upgrades
- Prevent access to newer OpenShift features that require crun
- Result in unsupported configurations in future releases

## Diagnosis

To identify which nodes are using runc, query the CRI-O metrics:

```bash
$ oc get nodes
$ oc debug node/<node-name>
# chroot /host
# crio status config | grep -A 2 "default_runtime"
```

Look for the `default_runtime` setting. If it shows `default_runtime = "runc"`,
the node is using runc.

To see all nodes using runc via metrics:

```bash
$ oc -n openshift-monitoring exec -it prometheus-k8s-0 -c prometheus -- sh
$ promtool query instant http://localhost:9090 \
  'container_runtime_crio_default_runtime{runtime="runc"}'
```

## Mitigation

To migrate from runc to crun, you can use a ContainerRuntimeConfig custom
resource (CR). Starting in OpenShift Container Platform 4.18, crun is the
default container runtime for new installations.

### For Clusters Upgraded from OpenShift Container Platform 4.17

If your cluster was upgraded from OpenShift Container Platform 4.17,
the runc container runtime remains unchanged as the default.
During the upgrade, two MachineConfig objects were created to
override the new default runtime. You can migrate to crun by
deleting these MachineConfig objects:

Migrate worker nodes to crun:
```bash
$ oc delete machineconfig \ 
  00-override-worker-generated-crio-default-container-runtime
```

Migrate control plane nodes to crun:
```bash
$ oc delete machineconfig \
  00-override-master-generated-crio-default-container-runtime
```

### For Any OpenShift Container Platform 4.18 or Greater Cluster

You can configure crun as the container runtime for specific nodes by creating
a ContainerRuntimeConfig CR.

#### Procedure

1. Create a YAML file for the ContainerRuntimeConfig CR for worker nodes:

```yaml
apiVersion: machineconfiguration.openshift.io/v1
kind: ContainerRuntimeConfig
metadata:
  name: configure-crun-worker
spec:
  machineConfigPoolSelector:
    matchLabels:
      pools.operator.machineconfiguration.openshift.io/worker: ''
  containerRuntimeConfig:
    defaultRuntime: "crun"
```

For control plane nodes, use a separate ContainerRuntimeConfig CR:

```yaml
apiVersion: machineconfiguration.openshift.io/v1
kind: ContainerRuntimeConfig
metadata:
  name: configure-crun-master
spec:
  machineConfigPoolSelector:
    matchLabels:
      pools.operator.machineconfiguration.openshift.io/master: ''
  containerRuntimeConfig:
    defaultRuntime: "crun"
```

2. Create the ContainerRuntimeConfig CR:

```bash
$ oc create -f configure-crun-worker.yaml
$ oc create -f configure-crun-master.yaml
```

3. Verify that the CR is created:

```bash
$ oc get ContainerRuntimeConfig
```

4. Check that new containerruntime machine configs are created:

```bash
$ oc get machineconfigs | grep containerrun
```

5. Monitor the machine config pools until all are shown as ready:

```bash
$ oc get mcp worker
$ oc get mcp master
```

Wait for all machine config pools to show `UPDATED=True` and `UPDATING=False`.

#### Verification

After the nodes return to a ready state, verify the changes:

1. Open an oc debug session to a node:

```bash
$ oc debug node/<node-name>
```

2. Set /host as the root directory within the debug shell:

```bash
sh-5.1# chroot /host
```

3. Check the container runtime:

```bash
sh-5.1# crio status config | grep default_runtime
```

Example output:
```bash
default_runtime = "crun"
```

The alert should automatically resolve once all nodes are using crun and the
metric `container_runtime_crio_default_runtime{runtime="runc"}` is no longer present.

### Additional Resources

For more detailed information on configuring the container runtime, see the
[OpenShift documentation on machine configuration](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/machine_configuration/machine-configs-custom#config-container-runtime_machine-configs-custom).
