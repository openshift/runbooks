# SystemMemoryExceedsReservation

## Meaning

This alert is triggered when a node is deteced to
be using more memory then is reserved by the
kubelet. This alert will fire if the node
is exceeding 95% of the reservation for 15
minutes or more.

## Impact

This alert is a warning to let a system admin
know the node is filling up all its allocatable
memory. The node needs this memory in order to
run and satisfy system processes. This alert
firing does not nessarily mean the node is
resource exhausted at the moment.

## Diagnosis

This alert is triggered by a expression
which consists of a Prometheus query. The
full query is as follows.

```console
sum by (node) (container_memory_rss{id="/system.slice"}) > ((sum by (node) (kube_node_status_capacity{resource="memory"} - kube_node_status_allocatable{resource="memory"})) * 0.95)
```

This can be split into two halfs. Using
the comparator operator, "greater then: `>`" as the split
point.

- The right side of the query consists of
static values.

  `kube_node_status_capacity{resource="memory"}`
  (The total memory capacity of the node)

  minus

  `kube_node_status_allocatable{resource="memory"}`
  (The amount of memory allocatable for pods)

  This gives the amount of memory allocated
to the nodes. This is then multiplied by 0.95
to get the 95th percentile.

- The left side of the query consists of a
dynamic value.

  `container_memory_rss{id="/system.slice"}`
(The total resident set slice which is a
portion of the system's memory occupied by
a process that is held in the main memory)

  if this value is greather then the 95th
  percentile of the allocatable memory for
  the node then the alert will go into pending.
  After 15 minutes in this state the alert
  will fire.

Since this is a comparator operation. If the
condition is not met, there will be no datapoints
displayed by the query.

## Mitigation

By default the `system-reserved` value
for memory is 1Gi. This value can be changed
manually post install. You can also have
the kubelet automatically determine and allocate the
system-reserved value via a script on each
node. This will take into account the CPU
and memory capacity that is installed on
the node.

To manually set the system-reserve value
or automatically set it. You must create a
KubeletConfig and give it the appropriate
`machineConfigPoolSelector` so it propagates
to the correct nodes you want to target.

- For manual allocation:
  
  ```console
  apiVersion: machineconfiguration.openshift.io/v1
  kind: KubeletConfig
  metadata:
    name: set-allocatable
  spec:
    machineConfigPoolSelector:
      matchLabels:
        pools.operator.machineconfiguration.openshift.io/worker: ""
    kubeletConfig:
      systemReserved:
        cpu: 1000m
        memory: 3Gi
  ```
- For automatic allocation:

  ```console
  apiVersion: machineconfiguration.openshift.io/v1
  kind: KubeletConfig
  metadata:
    name: dynamic-node
  spec:
    autoSizingReserved: true
    machineConfigPoolSelector:
      matchLabels:
        pools.operator.machineconfiguration.openshift.io/worker: ""
  ```
If increasing the memory value for
system-reserved is not an option.
Then you'll need to investigate
and troubleshoot which processes
are consuming the hosts memory.

Some useful commands for troubleshooting.

- You can use the `top` command on
the host to get a dynamic update of
the largest memory consuming proccesses.
For instance to get the top 100 memory
consuming processes on a node.

   ```console
     $ oc debug node/<node>
     $ chroot /host
     $ top -b -n100 -d2
   ```

- Another host level command is the `free`
command which allows you to check the memory
statistics of the node.

- Each node also contains a file called
`/proc/meminfo`. This file provides a usage
report about memory on the system. You can
learn how to interperet the fields [here](https://access.redhat.com/solutions/406773).

- For kubelet level commands you can get
the memory usage of individual pods by
using the `oc adm top pods` command.
You can further tune it to look at
individual containers by adding the
`--containers=true` flag.

- Prometheus itself can be used to deep
dive into the memory usage of nodes and
pods. Red Hat provides multiple pre-built
dashboards and PromQL queries to track
memory usage over time. All within the
`Observe` section of the OpenShift
console.