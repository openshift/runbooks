# NodeClockNotSynchronising

## Meaning

The `NodeClockNotSynchronising` alert triggers when a node is affected by
issues with the NTP server for that node. For example, this alert might trigger
when certificates are rotated for the API Server on a node, and the
certificates fail validation because of an invalid time.


## Impact
This alert is critical. It indicates an issue that can lead to the API Server
Operator becoming degraded or unavailable. If the API Server Operator becomes
degraded or unavailable, this issue can negatively affect other Operators, such
as the Cluster Monitoring Operator.

## Diagnosis

To diagnose the underlying issue, start a debug pod on the affected node and
check the `chronyd` service:

```shell
oc -n default debug node/<affected_node_name>
chroot /host
systemctl status chronyd
```
### Bug#1822211 

https://access.redhat.com/solutions/4976641
Issue
    When trying to run oc debug node over a tainted node with NoExecute effect, the debug pod got terminated.
    Error trying to debug a node when the node is tainted:
    error: unable to create the debug pod "worker-0.example.redhat.com-debug"
Resolution
    Add a toleration to a "dummy" namespace, and let the debug pod to be created in this namespace with --to-namespace option:
```shell
$ oc new-project dummy
$ oc patch namespace dummy --type=merge -p '{"metadata": {"annotations": { "scheduler.alpha.kubernetes.io/defaultTolerations": "[{\"operator\": \"Exists\"}]"}}}'
$ oc debug node/worker-0.example.redhat.com --to-namespace dummy

## Mitigation

1. If the `chronyd` service is failing or stopped, start it:

    ```shell
    systemctl start chronyd
    ```
    If the chronyd service is ready, restart it

    ```shell
    systemctl restart chronyd
    ```

    If `chronyd` starts or restarts successfuly, the service adjusts the clock
    and displays something similar to the following example output:

    ```shell
    Oct 18 19:39:36 ip-100-67-47-86 chronyd[2055318]: System clock wrong by 16422.107473 seconds, adjustment started
    Oct 19 00:13:18 ip-100-67-47-86 chronyd[2055318]: System clock was stepped by 16422.107473 seconds
    ```

2. Verify that the `chronyd` service is running:

    ```shell
    systemctl status chronyd
    ```

3. Verify using PromQL:

    ```console
    min_over_time(node_timex_sync_status[5m])
    node_timex_maxerror_seconds
    ```
    `node_timex_sync_status` returns `1` if NTP is working properly,or `0` if
    NTP is not working properly. `node_timex_maxerror_seconds` indicates how
    many seconds NTP is falling behind.

    The alert triggers when the value for
    `min_over_time(node_timex_sync_status[5m])` equals `0` and the value for
    `node_timex_maxerror_seconds` is greater than or equal to `16`.
