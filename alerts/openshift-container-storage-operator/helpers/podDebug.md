# Pod debugging

pod status: pending → Check for resource issues, pending pvcs, node assignment,
kubelet problems.

```bash
    oc project openshift-storage
    oc get pod | grep {ceph-component}
```

Set MYPOD for convenience:

```bash
    # Examine the output for a {ceph-component} that is in the pending state, not running or not ready
    MYPOD=<pod identified as the problem pod>
```

Look for resource limitations or pending pvcs. Otherwise, check for node
assignment.

```bash
    oc get pod/${MYPOD} -o wide
```

pod status: NOT pending, running, but NOT ready → Check readiness probe.

```bash
    oc describe pod/${MYPOD}
```

pod status: NOT pending, but NOT running → Check for app or image issues.

```bash
    oc logs pod/${MYPOD}
```

If a node was assigned, check kubelet on the node.