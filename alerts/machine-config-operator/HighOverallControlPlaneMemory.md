# HighOverallControlPlaneMemory

## Meaning

This alert becomes pending when the total memory
in use across all control plane nodes is detected
to be using higher than 60% of total available memory.
If the control plane nodes continue to use memory above the
threshold for 1 hour, the alert will fire.

## Impact
The memory usage per instance within control
plane nodes influences the stability
and responsiveness of the cluster, most noticably in the etcd and
Kubernetes API server pods. Moreover, OOM kill can occur
with excessive memory usage, which negatively
influences the pod scheduling. Etcd also relies on a certain number of
nodes being in quorum in order to operate.
If multiple control plane nodes become unstable,
quorum can be lost and the cluster can cease to operate.

## Diagnosis

This alert is triggered by an expression
that consists of a Prometheus query.

Targeting any node with the role of master,
prometheus computes the free memory left
on all nodes, the buffer memory on all nodes,
and the cached memory on all nodes. This is
divided by the total amount of memory available on all nodes.

Prometheus then multiplies the resulting value by 100
and this value is the total percentage
of memory used by the control plane nodes. If this value
is larger than the 60% for a certain period
of time, Prometheus will fire the alert.

To diagnose why a control plane node is using a high
amount of memory, you can investigate the following:

- The `oc describe` command will retrieve a large
variety of data on nodes, such as pods
scheduled on the node and the current
memory/cpu requests of those pods.
The command will also list the total
allocated resources on the node and
how much it is currently using. You can run
the following script to get a list of all control plane
nodes and to see only the memory and allocation
statistics.

  ```console
  $ oc get nodes -l node-role.kubernetes.io/control-plane -o name | while read NODE; do printf -- '\n----------%s description---------\n\n' "${NODE}"; oc describe "${NODE}"; done
  ```
- For kubelet-level commands, you can get
the memory usage of individual pods by
using the `oc adm top pods` command.
You can further tune the command to look at
individual containers by adding the
`--containers=true` flag and also
sort by memory from highest to
lowest by using the `--sort-by=memory` flag.

- For troubleshooting individual nodes,
you can SSH or `oc debug` into the control
plane node and use commands such as `top`
to diagnose memory usage.

  ```console
  $ oc debug node/<node>
  $ chroot /host
  $ top -b -n100 -d2
  ```
- You can use Prometheus to deep
dive into the memory usage of nodes and
pods. Red Hat provides multiple pre-built
dashboards and PromQL queries to track
memory usage over time. All within the
`Observe` section of the OpenShift
console.

- The memory use of `etcd` could be another
reason that memory could be high on
the control plane. Acting as the primary datastore
for the cluster state, etcd is
constantly being checked and updated by the Kube API to
match the current cluster state. An etcd pod runs on each control
plane node. Because of the amount of data
and querying etcd keeps track of,
this can lead to high memory usage.
To ensure etcd is running smoothly,
you can check the response time of etcd and
health status with commands run directly on the etcd pods.
  ```console
  $ oc -n openshift-etcd rsh <etcd-pod-name>
  # etcdctl member list -w table
  # etcdctl endpoint health -w table
  # etcdctl endpoint status -w table
  ```
  You can also get the amount of objects that etcd is storing by
  using this command.
  ```console
  $ oc exec -n openshift-etcd <etcd-pod> -c etcd -n openshift-etcd -- etcdctl --command-timeout=160s get / --prefix --keys-only |sed '/^$/d' | cut -d/ -f3,4 | sort | uniq -c | sort -rn |head
  ```
  You can also use the [fio tool](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/scalability_and_performance/recommended-performance-and-scalability-practices-2#etcd-verify-hardware_recommended-etcd-practices)
  to check the disk performance for the storage disks backing etcd.

## Mitigation
- Ensure your control plane nodes
are appropriately sized for your environment. Red Hat
offers [sizing and scaling guidelines](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/scalability_and_performance/recommended-performance-and-scalability-practices-2)
for clusters of various sizes.

- In the event this alert is firing
and the cluster it extremely unstable, you can
increase a control plane node's memory size
to allow the control plane to stabilize so you
can troubleshoot the issue.

- Because etcd is crucial to how the cluster runs, Red Hat
has documented [the recommended practices and requirements for optimal etcd performance.](https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/scalability_and_performance/recommended-performance-and-scalability-practices-2#recommended-etcd-practices_recommended-etcd-practices)

