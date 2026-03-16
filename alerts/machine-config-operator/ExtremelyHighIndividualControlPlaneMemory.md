# ExtremelyHighIndividualControlPlaneMemory

## Meaning
This alert becomes pending when a control plane
node is detected to be using more then 90%
of its total available memory. If the control
plane node continues to use memory above the
90% threshold for 45 minutes, the alert will fire.

## Impact
The memory usage per instance within control
plane nodes influences the stability
and responsiveness of the cluster, most visibily in the etcd and
Kubernetes API server pods. Moreover, OOM kill is expected
which negatively influences the pod scheduling.

## Diagnosis
This alert is triggered by an expression
that consists of a Prometheus query.

Targeting any node with the role of master,
prometheus sums the free memory left
on the node, the buffer memory on the node,
and the cached memory on the node. This is
divided by the total amount of memory available on that node.

Prometheus then multiplies the resulting value from each
master node by 100 and this value is the total percentage
of memory used by the control plane node. If this value
is larger than a certain percentage for a certain period
of time, Prometheus will fire the alert.

To diagnose why the control plane node is using a high
amount of memory, you can investigate the following:

- The `oc describe` command will retrieve a large
variety of data on nodes, such as pods
scheduled on the node and its current
memory/cpu requests. It will also list the total
allocated resources on the node and
how much it is currently using. You can run
a script to get a list of all control plane
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

- Another reason that the memory could
be high on the control planes is because of
`etcd`. Acting as the primary datastore for the cluster state, etcd is
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
  You can also get the amount of objects that etcd is storing via this command.
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
increase the control plane nodes memory size
to allow the control plane to stabilize so you
can troubleshoot the issue.

- Because etcd is crucial to how the cluster runs, Red Hat
has [documented](https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/scalability_and_performance/recommended-performance-and-scalability-practices-2#recommended-etcd-practices_recommended-etcd-practices)
the recommended practices and requirements for optimal etcd performance.
